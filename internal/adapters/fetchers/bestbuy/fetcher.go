package bestbuy

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"

    "golang.org/x/sync/errgroup"

    "PriceMonCA/internal/adapters/httpx"
    "PriceMonCA/internal/domain"
)

type Fetcher struct {
    cfg  Config
    http *httpx.Client
    log  *slog.Logger
}

func New(cfg Config, hc *httpx.Client, log *slog.Logger) (*Fetcher, error) {
    if hc == nil {
        return nil, fmt.Errorf("%w: http client required", domain.ErrInvalidInput)
    }
    if cfg.Host == "" || cfg.ProductPath == "" || cfg.LocationsPath == "" || cfg.AvailabilityPath == "" {
        return nil, fmt.Errorf("%w: bestbuy config incomplete", domain.ErrInvalidInput)
    }
    return &Fetcher{
        cfg:  cfg.withDefaults(),
        http: hc,
        log:  log.With("site", "bestbuy"),
    }, nil
}

func (f *Fetcher) Site() string { return "bestbuy" }

func (f *Fetcher) Fetch(ctx context.Context, ref domain.ProductRef) (domain.Product, error) {
    postal := ref.Hints["postalCodeFSA"]

    var (
        prod productDTO
        locs locationsDTO
    )

    // Independent → fetch concurrently.
    g, gctx := errgroup.WithContext(ctx)
    g.Go(func() error {
        u, err := f.productURL(ref.SKU)
        if err != nil {
            return err
        }
        return f.getJSON(gctx, u, &prod)
    })
    g.Go(func() error {
        u, err := f.locationsURL(postal)
        if err != nil {
            return err
        }
        return f.getJSON(gctx, u, &locs)
    })
    if err := g.Wait(); err != nil {
        return domain.Product{}, fmt.Errorf("bestbuy: fetch product/locations: %w", err)
    }

    // Availability depends on location IDs, so it must be sequential.
    availURL, err := f.availabilityURL([]string{ref.SKU}, postal, locs.StoreIDs())
    if err != nil {
        return domain.Product{}, err
    }
    var avail availabilityDTO
    if err := f.getJSON(ctx, availURL, &avail); err != nil {
        return domain.Product{}, fmt.Errorf("bestbuy: availability: %w", err)
    }

    return toDomain(ref, prod, locs, avail)
}

func (f *Fetcher) getJSON(ctx context.Context, url string, out any) error {
    body, err := f.http.GetJSON(ctx, url, map[string]string{
        "Accept-Language": f.cfg.Language,
    })
    if err != nil {
        return err
    }
    if err := json.Unmarshal(body, out); err != nil {
        return fmt.Errorf("decode %s: %w", url, err)
    }
    return nil
}
