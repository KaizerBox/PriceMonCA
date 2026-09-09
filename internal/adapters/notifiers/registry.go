package notifiers

import (
    "context"
    "fmt"
    "sync"

    "PriceMonCA/internal/domain"
)

type Notifier interface {
    Channel() string
    Notify(ctx context.Context, target domain.NotificationTarget, e domain.Event) error
}

type Registry struct {
    mu sync.RWMutex
    m  map[string]Notifier
}

func NewRegistry() *Registry { return &Registry{m: map[string]Notifier{}} }

func (r *Registry) Register(n Notifier) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.m[n.Channel()] = n
}

func (r *Registry) Get(channel string) (Notifier, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    n, ok := r.m[channel]
    if !ok {
        return nil, fmt.Errorf("%w: channel %q", domain.ErrNotFound, channel)
    }
    return n, nil
}
