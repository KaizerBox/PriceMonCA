package config

import (
    "fmt"
    "os"
    "time"

    "gopkg.in/yaml.v3"

    "PriceMonCA/internal/adapters/fetchers/bestbuy"
    "PriceMonCA/internal/adapters/notifiers/gmail"
)

type Config struct {
    LogLevel  string                `yaml:"logLevel"`
    HTTP      HTTPConfig            `yaml:"http"`
    Scheduler SchedulerConfig       `yaml:"scheduler"`
    Sites     map[string]SiteConfig `yaml:"sites"`
    Notifiers NotifiersConfig       `yaml:"notifiers"`
}

type HTTPConfig struct {
    ListenAddr string        `yaml:"listenAddr"`
    Timeout    time.Duration `yaml:"timeout"`
    RPS        float64       `yaml:"rps"`
    Burst      int           `yaml:"burst"`
    UserAgent  string        `yaml:"userAgent"`
    MaxTries   int           `yaml:"maxTries"`
}

type SchedulerConfig struct {
    DefaultSpec string `yaml:"defaultSpec"`
}

type SiteConfig struct {
    BestBuy bestbuy.Config `yaml:"bestbuy,omitempty"`
}

type NotifiersConfig struct {
    Gmail *gmail.Config `yaml:"gmail,omitempty"`
}

func Load() (Config, error) {
    path := os.Getenv("PRICEMON_CONFIG")
    if path == "" {
        path = "configs/pricemon.yaml"
    }
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("config: read %s: %w", path, err)
    }
    var c Config
    if err := yaml.Unmarshal(data, &c); err != nil {
        return Config{}, fmt.Errorf("config: parse: %w", err)
    }
    c.applyDefaults()
    return c, nil
}

func (c *Config) applyDefaults() {
    if c.LogLevel == "" {
        c.LogLevel = "info"
    }
    if c.HTTP.ListenAddr == "" {
        c.HTTP.ListenAddr = ":8080"
    }
    if c.HTTP.Timeout == 0 {
        c.HTTP.Timeout = 15 * time.Second
    }
    if c.Scheduler.DefaultSpec == "" {
        c.Scheduler.DefaultSpec = "@every 10m"
    }
}
