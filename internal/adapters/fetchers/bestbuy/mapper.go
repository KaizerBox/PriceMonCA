package bestbuy

import (
    "fmt"
    "strings"
    "time"

    "PriceMonCA/internal/domain"
)

const bestbuyDateFmt = "2006-01-02"
//Maps the Product, Store, and Availability API into DTO
func toDomain(ref domain.ProductRef, prod productDTO, locs locationsDTO, avail availabilityDTO) (domain.Product, error) {
    availBySKU := avail.BySKU()
    pa, ok := availBySKU[ref.SKU]
    if !ok {
        return domain.Product{}, fmt.Errorf("%w: no availability for sku %s", domain.ErrNotFound, ref.SKU)
    }

    storesByID := locs.IndexByID()
    inStore := make([]domain.StoreAvailability, 0, len(pa.Pickup.Details))
    for _, d := range pa.Pickup.Details {
        si, ok := storesByID[d.StoreID]
        if !ok {
            // Skip unknown stores instead of aborting the whole fetch.
            continue
        }
        inStore = append(inStore, domain.StoreAvailability{
            StoreID:        d.StoreID,
            StoreName:      si.StoreName,
            StoreAddress:   si.StoreAddress,
            QuantityOnHand: d.QuantityOnHand,
            HasInventory:   d.HasInventory,
            Reservable:     d.IsReservable,
            Fulfillable:    d.SupportsFulfillment,
        })
    }

    // Compute Available from real signals (this fixes the always-false bug).
    online := domain.OnlineAvailability{
        Status:            pa.Shipping.Status,
        QuantityRemaining: pa.Shipping.QuantityRemaining,
        Backorderable:     pa.Shipping.IsBackorderable,
        Available: strings.EqualFold(pa.Shipping.Status, "InStock") ||
            pa.Shipping.QuantityRemaining > 0 ||
            pa.Shipping.IsBackorderable,
    }
    if len(pa.Shipping.Details) > 0 {
        if t, err := time.Parse(bestbuyDateFmt, pa.Shipping.Details[0].DeliveryDate); err == nil {
            online.EstimatedDelivery = t
        }
    }

    currency := prod.Currency
    if currency == "" {
        currency = "CAD"
    }

    return domain.Product{
        Ref:          ref,
        Name:         prod.Name,
        RegularPrice: prod.RegularPrice,
        SalePrice:    prod.SalePrice,
        Currency:     currency,
        Seller:       firstNonEmpty(prod.SellerID, pa.SellerID),
        Online:       online,
        InStore:      inStore,
        FetchedAt:    time.Now().UTC(),
        SourceURL:    prod.ProductURL,
    }, nil
}

func firstNonEmpty(a, b string) string {
    if a != "" {
        return a
    }
    return b
}
