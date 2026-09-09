package bestbuy

type Config struct {
    Scheme           string            `yaml:"scheme"`
    Host             string            `yaml:"host"`
    ProductPath      string            `yaml:"productPath"`
    LocationsPath    string            `yaml:"locationsPath"`
    AvailabilityPath string            `yaml:"availabilityPath"`
    AvailabilityQ    map[string]string `yaml:"availabilityQuery"`
    Language         string            `yaml:"language"`
}

func (c Config) withDefaults() Config {
    if c.Scheme == "" {
        c.Scheme = "https"
    }
    if c.Language == "" {
        c.Language = "en-CA"
    }
    if c.AvailabilityQ == nil {
        c.AvailabilityQ = map[string]string{}
    }
    return c
}
