package domain

import "time"

//Notification subscription
type Subscription struct {
    ID        string               `json:"id"`
    UserID    string               `json:"userId"`
    Ref       ProductRef           `json:"ref"`
    Rule      Rule                 `json:"rule"`
    Targets   []NotificationTarget `json:"targets"`
    CronSpec  string               `json:"cronSpec,omitempty"` // empty = scheduler default
    Active    bool                 `json:"active"`
    CreatedAt time.Time            `json:"createdAt"`
}

type NotificationTarget struct {
    Channel string `json:"channel"` // "email", "discord", ...
    Address string `json:"address"` // email address / webhook URL / etc.
}
