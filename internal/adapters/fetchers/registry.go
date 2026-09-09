package fetchers

import (
    "context"
    "fmt"
    "sync"

    "pricemon/internal/domain"
)

// Fetcher is the interface every site adapter implements.
type Fetcher interface {
    Site() string
    Fetch(ctx context.Context, ref domain.ProductRef) (domain.Product, error)
}

// Registry maps site key -> Fetcher. Concrete, thread-safe, no interface needed.
type Registry struct {
    mu sync.RWMutex
    m  map[string]Fetcher
}

func NewRegistry() *Registry { return &Registry{m: map[string]Fetcher{}} }

func (r *Registry) Register(f Fetcher) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.m[f.Site()] = f
}

func (r *Registry) Get(site string) (Fetcher, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    f, ok := r.m[site]
    if !ok {
        return nil, fmt.Errorf("%w: site %q", domain.ErrNotFound, site)
    }
    return f, nil
}
