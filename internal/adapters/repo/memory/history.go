package memory

import (
    "context"
    "fmt"
    "sync"

    "PriceMonCA/internal/domain"
)

type PriceHistoryRepo struct {
    mu sync.RWMutex
    m  map[string][]domain.Product
}

func NewPriceHistoryRepo() *PriceHistoryRepo {
    return &PriceHistoryRepo{m: map[string][]domain.Product{}}
}

func key(ref domain.ProductRef) string { return ref.Site + "|" + ref.SKU }

func (r *PriceHistoryRepo) Append(_ context.Context, p domain.Product) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    k := key(p.Ref)
    r.m[k] = append(r.m[k], p)
    return nil
}

func (r *PriceHistoryRepo) Latest(_ context.Context, ref domain.ProductRef) (domain.Product, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    xs := r.m[key(ref)]
    if len(xs) == 0 {
        return domain.Product{}, fmt.Errorf("%w: history for %s/%s", domain.ErrNotFound, ref.Site, ref.SKU)
    }
    return xs[len(xs)-1], nil
}
