package domain

import "time"

// ProductRef identifies a product on a specific site.
// Hints hold site-specific parameters (postal code, region, etc.).
type ProductRef struct {
    Site  string            `json:"site"`
    SKU   string            `json:"sku"`
    Hints map[string]string `json:"hints,omitempty"`
}

// Product is the normalized view of a product across all sites.
type Product struct {
    Ref          ProductRef          `json:"ref"`
    Name         string              `json:"name"`
    RegularPrice float64             `json:"regularPrice"`
    SalePrice    float64             `json:"salePrice"`
    Currency     string              `json:"currency"`
    Seller       string              `json:"seller"`
    Online       OnlineAvailability  `json:"online"`
    InStore      []StoreAvailability `json:"inStore,omitempty"`
    FetchedAt    time.Time           `json:"fetchedAt"`
    SourceURL    string              `json:"sourceUrl,omitempty"`
}

func (p Product) Discount() float64 { return p.RegularPrice - p.SalePrice }

func (p Product) DiscountPct() float64 {
    if p.RegularPrice <= 0 {
        return 0
    }
    return p.Discount() / p.RegularPrice * 100
}
//Online availability info of the the product
type OnlineAvailability struct {
    Available         bool      `json:"available"`
    Status            string    `json:"status"`
    QuantityRemaining int       `json:"quantityRemaining"`
    EstimatedDelivery time.Time `json:"estimatedDelivery,omitempty"`
    Backorderable     bool      `json:"backorderable"`
}
////In Store availability info of the the product
type StoreAvailability struct {
    StoreID        string `json:"storeId"`
    StoreName      string `json:"storeName"`
    StoreAddress   string `json:"storeAddress"`
    QuantityOnHand int    `json:"quantityOnHand"`
    HasInventory   bool   `json:"hasInventory"`
    Reservable     bool   `json:"reservable"`
    Fulfillable    bool   `json:"fulfillable"`
}
