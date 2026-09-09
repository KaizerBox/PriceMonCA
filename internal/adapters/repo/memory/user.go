package memory

import (
    "context"
    "fmt"
    "sync"

    "pricemon/internal/domain"
)

type UserRepo struct {
    mu sync.RWMutex
    m  map[string]domain.User
}

func NewUserRepo() *UserRepo { return &UserRepo{m: map[string]domain.User{}} }

func (r *UserRepo) Put(_ context.Context, u domain.User) error {
    if u.ID == "" {
        return fmt.Errorf("%w: user id required", domain.ErrInvalidInput)
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    r.m[u.ID] = u
    return nil
}

func (r *UserRepo) Get(_ context.Context, id string) (domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    u, ok := r.m[id]
    if !ok {
        return domain.User{}, fmt.Errorf("%w: user %s", domain.ErrNotFound, id)
    }
    return u, nil
}
