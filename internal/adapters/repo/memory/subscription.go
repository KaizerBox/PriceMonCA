package memory

import (
    "context"
    "fmt"
    "sync"

    "pricemon/internal/domain"
)

type SubscriptionRepo struct {
    mu sync.RWMutex
    m  map[string]domain.Subscription
}

func NewSubscriptionRepo() *SubscriptionRepo {
    return &SubscriptionRepo{m: map[string]domain.Subscription{}}
}

func (r *SubscriptionRepo) Put(_ context.Context, s domain.Subscription) error {
    if s.ID == "" {
        return fmt.Errorf("%w: subscription id required", domain.ErrInvalidInput)
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    r.m[s.ID] = s
    return nil
}

func (r *SubscriptionRepo) Get(_ context.Context, id string) (domain.Subscription, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    s, ok := r.m[id]
    if !ok {
        return domain.Subscription{}, fmt.Errorf("%w: subscription %s", domain.ErrNotFound, id)
    }
    return s, nil
}

func (r *SubscriptionRepo) ListActive(_ context.Context) ([]domain.Subscription, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    out := make([]domain.Subscription, 0, len(r.m))
    for _, s := range r.m {
        if s.Active {
            out = append(out, s)
        }
    }
    return out, nil
}

func (r *SubscriptionRepo) Delete(_ context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.m[id]; !ok {
        return fmt.Errorf("%w: subscription %s", domain.ErrNotFound, id)
    }
    delete(r.m, id)
    return nil
}
