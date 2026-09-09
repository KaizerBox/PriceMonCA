package domain

type RuleKind string

const (
    RulePrice    RuleKind = "price"
    RuleStock    RuleKind = "stock"
    RuleDiscount RuleKind = "discount"
    RuleAll      RuleKind = "all"
    RuleAny      RuleKind = "any"
)

type Rule struct {
    ID   string   `json:"id"`
    Kind RuleKind `json:"kind"`

    Price    *PriceRule    `json:"price,omitempty"`
    Stock    *StockRule    `json:"stock,omitempty"`
    Discount *DiscountRule `json:"discount,omitempty"`
    All      []Rule        `json:"all,omitempty"`
    Any      []Rule        `json:"any,omitempty"`
}

// PriceRule matches when SalePrice is within [Min, Max]. Nil bound = open.
type PriceRule struct {
    Min *float64 `json:"min,omitempty"`
    Max *float64 `json:"max,omitempty"`
}

// StockRule matches on availability and quantity bounds.
type StockRule struct {
    InStock *bool `json:"inStock,omitempty"`
    Min     *int  `json:"min,omitempty"`
    Max     *int  `json:"max,omitempty"`
}

// DiscountRule matches on the percentage discount.
type DiscountRule struct {
    MinPct *float64 `json:"minPct,omitempty"`
    MaxPct *float64 `json:"maxPct,omitempty"`
}

// Decision is the evaluator's verdict for a single rule against a product.
type Decision struct {
    Match   bool     `json:"match"`
    Reason  string   `json:"reason"`
    Matched []string `json:"matched,omitempty"`
}
