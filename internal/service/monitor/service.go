package monitor

import (
    "context"
    "errors"
    "log/slog"

    "PriceMonCA/internal/adapters/fetchers"
    "PriceMonCA/internal/adapters/notifiers"
    "PriceMonCA/internal/domain"
)

// Consumer-side interfaces (small, local).

type FetcherRegistry interface {
    Get(site string) (fetchers.Fetcher, error)
}

type NotifierRegistry interface {
    Get(channel string) (notifiers.Notifier, error)
}

type Evaluator interface {
    Evaluate(ctx context.Context, p domain.Product, r domain.Rule) (domain.Decision, error)
}

type PriceHistoryRepo interface {
    Append(ctx context.Context, p domain.Product) error
}

// Service is concrete — one implementation only.
type Service struct {
    fetchers  FetcherRegistry
    eval      Evaluator
    notifiers NotifierRegistry
    history   PriceHistoryRepo
    log       *slog.Logger
}

func New(f FetcherRegistry, e Evaluator, n NotifierRegistry, h PriceHistoryRepo, log *slog.Logger) *Service {
    return &Service{
        fetchers:  f,
        eval:      e,
        notifiers: n,
        history:   h,
        log:       log.With("component", "monitor"),
    }
}

// RunOne executes one full cycle: fetch → persist → evaluate → notify.
func (s *Service) RunOne(ctx context.Context, sub domain.Subscription) error {
    log := s.log.With("sub", sub.ID, "site", sub.Ref.Site, "sku", sub.Ref.SKU)

    f, err := s.fetchers.Get(sub.Ref.Site)
    if err != nil {
        return err
    }

    prod, err := f.Fetch(ctx, sub.Ref)
    if err != nil {
        log.Error("fetch failed", "err", err)
        return err
    }

    if err := s.history.Append(ctx, prod); err != nil {
        log.Warn("history append failed", "err", err) // best-effort
    }

    dec, err := s.eval.Evaluate(ctx, prod, sub.Rule)
    if err != nil {
        log.Error("evaluate failed", "err", err)
        return err
    }
    if !dec.Match {
        log.Debug("no match", "reason", dec.Reason)
        return nil
    }

    ev := domain.EventFromMatch(prod, sub, dec)
    var errs []error
    for _, t := range sub.Targets {
        n, err := s.notifiers.Get(t.Channel)
        if err != nil {
            log.Warn("notifier not registered", "channel", t.Channel)
            continue
        }
        if err := n.Notify(ctx, t, ev); err != nil {
            errs = append(errs, err)
            log.Error("notify failed", "channel", t.Channel, "err", err)
        }
    }
    return errors.Join(errs...)
}
