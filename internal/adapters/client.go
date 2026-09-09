package httpx

import (
    "context"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"

    "golang.org/x/time/rate"
)

type Client struct {
    hc       *http.Client
    limiter  *rate.Limiter
    ua       string
    maxTries int
}

type Options struct {
    Timeout   time.Duration
    RPS       float64
    Burst     int
    UserAgent string
    MaxTries  int
}

func New(opts Options) *Client {
    if opts.Timeout == 0 {
        opts.Timeout = 15 * time.Second
    }
    if opts.RPS == 0 {
        opts.RPS = 2
    }
    if opts.Burst == 0 {
        opts.Burst = 4
    }
    if opts.MaxTries <= 0 {
        opts.MaxTries = 3
    }
    if opts.UserAgent == "" {
        opts.UserAgent = "pricemon/1.0"
    }
    return &Client{
        hc:       &http.Client{Timeout: opts.Timeout},
        limiter:  rate.NewLimiter(rate.Limit(opts.RPS), opts.Burst),
        ua:       opts.UserAgent,
        maxTries: opts.MaxTries,
    }
}

// StatusError is a typed error for non-2xx responses.
type StatusError struct {
    StatusCode int
    Body       []byte
}

func (e *StatusError) Error() string {
    return fmt.Sprintf("unexpected status %d: %s", e.StatusCode, truncate(e.Body, 256))
}

func (e *StatusError) Retryable() bool {
    return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

func (c *Client) GetJSON(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
    var lastErr error
    for attempt := 0; attempt < c.maxTries; attempt++ {
        if err := c.limiter.Wait(ctx); err != nil {
            return nil, err
        }

        req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
        if err != nil {
            return nil, err
        }
        req.Header.Set("User-Agent", c.ua)
        req.Header.Set("Accept", "application/json")
        for k, v := range headers {
            req.Header.Set(k, v)
        }

        resp, err := c.hc.Do(req)
        if err != nil {
            lastErr = err
            if ctx.Err() != nil {
                return nil, ctx.Err()
            }
            backoff(attempt)
            continue
        }
        body, readErr := io.ReadAll(resp.Body)
        _ = resp.Body.Close()
        if readErr != nil {
            lastErr = readErr
            backoff(attempt)
            continue
        }

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            return body, nil
        }
        se := &StatusError{StatusCode: resp.StatusCode, Body: body}
        lastErr = se
        if !se.Retryable() {
            return nil, se
        }
        backoff(attempt)
    }
    return nil, errors.Join(errors.New("giving up after retries"), lastErr)
}

func backoff(attempt int) {
    d := time.Duration(1<<attempt) * 200 * time.Millisecond
    time.Sleep(d)
}

func truncate(b []byte, n int) []byte {
    if len(b) <= n {
        return b
    }
    return b[:n]
}
