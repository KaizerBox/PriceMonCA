package scheduler

import (
    "context"
    "fmt"
    "log/slog"
    "sync"
    "time"

    "github.com/robfig/cron/v3"

    "PriceMonCA/internal/domain"
)

// SubscriptionLister is defined here (consumer side) — only the methods we need.
type SubscriptionLister interface {
    ListActive(ctx context.Context) ([]domain.Subscription, error)
    Get(ctx context.Context, id string) (domain.Subscription, error)
}

// Runner is what the scheduler invokes per tick per subscription.
type Runner interface {
    RunOne(ctx context.Context, sub domain.Subscription) error
}

const defaultJobTimeout = 2 * time.Minute

type Scheduler struct {
    cron        *cron.Cron
    log         *slog.Logger
    subs        SubscriptionLister
    runner      Runner
    defaultSpec string

    mu      sync.Mutex
    entries map[string]cron.EntryID
    running bool
}

func New(log *slog.Logger, subs SubscriptionLister, runner Runner, defaultSpec string) *Scheduler {
    if defaultSpec == "" {
        defaultSpec = "@every 10m"
    }
    return &Scheduler{
        cron:        cron.New(),
        log:         log.With("component", "scheduler"),
        subs:        subs,
        runner:      runner,
        defaultSpec: defaultSpec,
        entries:     map[string]cron.EntryID{},
    }
}

// Start loads all active subscriptions and schedules them.
func (s *Scheduler) Start(ctx context.Context) error {
    s.mu.Lock()
    if s.running {
        s.mu.Unlock()
        return nil
    }
    s.running = true
    s.mu.Unlock()

    list, err := s.subs.ListActive(ctx)
    if err != nil {
        return fmt.Errorf("scheduler: load subscriptions: %w", err)
    }
    for _, sub := range list {
        if err := s.ScheduleSubscription(sub); err != nil {
            s.log.Warn("could not schedule", "sub", sub.ID, "err", err)
        }
    }
    s.cron.Start()
    s.log.Info("scheduler started", "count", len(list))
    return nil
}

// ScheduleSubscription registers or replaces a subscription's job (idempotent).
func (s *Scheduler) ScheduleSubscription(sub domain.Subscription) error {
    spec := sub.CronSpec
    if spec == "" {
        spec = s.defaultSpec
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    if id, ok := s.entries[sub.ID]; ok {
        s.cron.Remove(id)
        delete(s.entries, sub.ID)
    }

    id, err := s.cron.AddFunc(spec, func() {
        jobCtx, cancel := context.WithTimeout(context.Background(), defaultJobTimeout)
        defer cancel()
        if err := s.runner.RunOne(jobCtx, sub); err != nil {
            s.log.Error("job failed", "sub", sub.ID, "err", err)
        }
    })
    if err != nil {
        return fmt.Errorf("scheduler: add cron %q: %w", spec, err)
    }
    s.entries[sub.ID] = id
    s.log.Info("scheduled", "sub", sub.ID, "spec", spec)
    return nil
}

func (s *Scheduler) UnscheduleSubscription(id string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if eid, ok := s.entries[id]; ok {
        s.cron.Remove(eid)
        delete(s.entries, id)
    }
}

func (s *Scheduler) Stop() {
    <-s.cron.Stop().Done()
    s.log.Info("scheduler stopped")
}
