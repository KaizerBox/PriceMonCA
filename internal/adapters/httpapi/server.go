package httpapi

import (
    "log/slog"
    "net/http"
    "time"

    "PriceMonCA/internal/adapters/repo/memory"
    "PriceMonCA/internal/adapters/scheduler"
)

type Server struct {
    log   *slog.Logger
    subs  *memory.SubscriptionRepo
    users *memory.UserRepo
    sched *scheduler.Scheduler
}

func NewServer(log *slog.Logger, subs *memory.SubscriptionRepo, users *memory.UserRepo, sched *scheduler.Scheduler) *Server {
    return &Server{log: log, subs: subs, users: users, sched: sched}
}

func (s *Server) Handler() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    mux.HandleFunc("POST /v1/subscriptions", s.createSubscription)
    mux.HandleFunc("GET /v1/subscriptions/{id}", s.getSubscription)
    mux.HandleFunc("DELETE /v1/subscriptions/{id}", s.deleteSubscription)

    return withCommon(s.log, mux)
}

func withCommon(log *slog.Logger, h http.Handler) http.Handler {
    return recoverer(log)(requestLogger(log)(timeout(30*time.Second)(h)))
}
