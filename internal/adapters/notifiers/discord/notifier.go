package discord

// TODO:
// - Config { WebhookURL string }
// - Channel() returns "discord"
// - Notify(): POST JSON {"content": e.Body} to target.Address (webhook URL)
//   using the shared *httpx.Client (retries + rate-limit for free).
// - If target.Address is empty, fall back to cfg.WebhookURL.
