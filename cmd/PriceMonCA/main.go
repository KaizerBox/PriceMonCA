package main

import (
    "context"
    "errors"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "PriceMonCA/internal/adapters/api/httpapi"
    "PriceMonCA/internal/adapters/fetchers"
    "PriceMonCA/internal/adapters/fetchers/bestbuy"
    "PriceMonCA/internal/adapters/httpx"
    "PriceMonCA/internal/adapters/notifiers"
    "PriceMonCA/internal/adapters/notifiers/gmail"
    "PriceMonCA/internal/adapters/repo/memory"
    "PriceMonCA/internal/adapters/rules"
    "PriceMonCA/internal/adapters/scheduler"
    "PriceMonCA/internal/config"
    "PriceMonCA/internal/obs"
    "PriceMonCA/internal/service/monitor"
)

func main() {
    if err := run(); err != nil {
        slog.Error("fatal", "err", err)
        os.Exit(1)
    }
}

func run() error {
    cfg, err := config.Load()
    if err != nil {
        return err
    }

    log := obs.NewLogger(cfg.LogLevel)

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    httpClient := httpx.New(httpx.Options{
        Timeout:   cfg.HTTP.Timeout,
        RPS:       cfg.HTTP.RPS,
        Burst:     cfg.HTTP.Burst,
        UserAgent: cfg.HTTP.UserAgent,
        MaxTries:  cfg.HTTP.MaxTries,
    })

    fetcherReg := fetchers.NewRegistry()
    for site, sc := range cfg.Sites {
        switch site {
        case "bestbuy":
            f, err := bestbuy.New(sc.BestBuy, httpClient, log)
            if err != nil {
                return err
            }
            fetcherReg.Register(f)
        default:
            log.Warn("unknown site in config, skipping", "site", site)
        }
    }

    notifierReg := notifiers.NewRegistry()
    if cfg.Notifiers.Gmail != nil {
        n, err := gmail.New(ctx, *cfg.Notifiers.Gmail, log)
        if err != nil {
            return err
        }
        notifierReg.Register(n)
    }

    subRepo := memory.NewSubscriptionRepo()
    histRepo := memory.NewPriceHistoryRepo()
    userRepo := memory.NewUserRepo()

    eval := rules.NewEngine()
    svc := monitor.New(fetcherReg, eval, notifierReg, histRepo, log)

    sched := scheduler.New(log, subRepo, svc, cfg.Scheduler.DefaultSpec)
    if err := sched.Start(ctx); err != nil {
        return err
    }
    defer sched.Stop()

    srv := &http.Server{
        Addr:              cfg.HTTP.ListenAddr,
        Handler:           httpapi.NewServer(log, subRepo, userRepo, sched).Handler(),
        ReadHeaderTimeout: 5 * time.Second,
    }

    errCh := make(chan error, 1)
    go func() {
        log.Info("http listening", "addr", cfg.HTTP.ListenAddr)
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            errCh <- err
        }
    }()

    select {
    case <-ctx.Done():
        log.Info("shutdown signal received")
    case err := <-errCh:
        log.Error("http server error", "err", err)
    }

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return srv.Shutdown(shutdownCtx)
}
