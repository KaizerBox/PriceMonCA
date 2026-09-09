package httpapi

import (
    "log/slog"
    "net/http"
    "runtime/debug"
    "time"
)

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
            next.ServeHTTP(sw, r)
            log.Info("request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", sw.status,
                "duration_ms", time.Since(start).Milliseconds(),
            )
        })
    }
}

func recoverer(log *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    log.Error("panic", "recover", rec, "stack", string(debug.Stack()))
                    http.Error(w, "internal error", http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}

func timeout(d time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.TimeoutHandler(next, d, "request timeout")
    }
}

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (s *statusWriter) WriteHeader(code int) {
    s.status = code
    s.ResponseWriter.WriteHeader(code)
}
