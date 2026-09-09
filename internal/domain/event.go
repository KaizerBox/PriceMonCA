package domain

import (
    "fmt"
    "time"
)
//Event used for notification
type Event struct {
    SubscriptionID string    `json:"subscriptionId"`
    Product        Product   `json:"product"`
    Decision       Decision  `json:"decision"`
    Subject        string    `json:"subject"`
    Body           string    `json:"body"`
    CreatedAt      time.Time `json:"createdAt"`
}

//Produce an event based on the match criteria
func EventFromMatch(p Product, sub Subscription, dec Decision) Event {
    subject := fmt.Sprintf("[pricemon] %s — %s", p.Name, dec.Reason)
    body := fmt.Sprintf(
        "Product: %s (%s / %s)\nPrice: %.2f %s (was %.2f, %.1f%% off)\nAvailable online: %v (qty %d)\nSource: %s\nReason: %s\n",
        p.Name, p.Ref.Site, p.Ref.SKU,
        p.SalePrice, p.Currency, p.RegularPrice, p.DiscountPct(),
        p.Online.Available, p.Online.QuantityRemaining,
        p.SourceURL, dec.Reason,
    )
    return Event{
        SubscriptionID: sub.ID,
        Product:        p,
        Decision:       dec,
        Subject:        subject,
        Body:           body,
        CreatedAt:      time.Now().UTC(),
    }
}
