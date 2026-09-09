package rules

import (
    "context"

    "PriceMonCA/internal/domain"
)

type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

func (Engine) Evaluate(_ context.Context, p domain.Product, r domain.Rule) (domain.Decision, error) {
    ok, reason, matched := evaluate(p, r)
    return domain.Decision{Match: ok, Reason: reason, Matched: matched}, nil
}

func evaluate(p domain.Product, r domain.Rule) (bool, string, []string) {
    switch r.Kind {
    case domain.RulePrice:
        if r.Price == nil {
            return false, "price rule missing body", nil
        }
        if r.Price.Min != nil && p.SalePrice < *r.Price.Min {
            return false, "price below min", nil
        }
        if r.Price.Max != nil && p.SalePrice > *r.Price.Max {
            return false, "price above max", nil
        }
        return true, "price in range", ids(r.ID)

    case domain.RuleStock:
        if r.Stock == nil {
            return false, "stock rule missing body", nil
        }
        if r.Stock.InStock != nil && *r.Stock.InStock != p.Online.Available {
            return false, "stock availability mismatch", nil
        }
        if r.Stock.Min != nil && p.Online.QuantityRemaining < *r.Stock.Min {
            return false, "quantity below min", nil
        }
        if r.Stock.Max != nil && p.Online.QuantityRemaining > *r.Stock.Max {
            return false, "quantity above max", nil
        }
        return true, "stock ok", ids(r.ID)

    case domain.RuleDiscount:
        if r.Discount == nil {
            return false, "discount rule missing body", nil
        }
        pct := p.DiscountPct()
        if r.Discount.MinPct != nil && pct < *r.Discount.MinPct {
            return false, "discount below min", nil
        }
        if r.Discount.MaxPct != nil && pct > *r.Discount.MaxPct {
            return false, "discount above max", nil
        }
        return true, "discount ok", ids(r.ID)

    case domain.RuleAll:
        var matched []string
        for _, sub := range r.All {
            ok, reason, m := evaluate(p, sub)
            if !ok {
                return false, "AND failed: " + reason, nil
            }
            matched = append(matched, m...)
        }
        return true, "all matched", matched

    case domain.RuleAny:
        for _, sub := range r.Any {
            if ok, _, m := evaluate(p, sub); ok {
                return true, "any matched", m
            }
        }
        return false, "no branch matched", nil
    }
    return false, "unknown rule kind", nil
}

func ids(s string) []string {
    if s == "" {
        return nil
    }
    return []string{s}
}
