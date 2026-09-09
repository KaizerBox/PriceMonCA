package bestbuy

import (
    "fmt"
    "net/url"
    "regexp"
    "strconv"
    "strings"

    "PriceMonCA/internal/domain"
)

// Canadian FSA: letter, digit, letter (with restricted first letter set).
var fsaRegex = regexp.MustCompile(`^[ABCEGHJKLMNPRSTVXY][0-9][A-Z]$`)
//Validate FSA of the Postal Code
func validateFSA(fsa string) error {
    if len(fsa) != 3 {
        return fmt.Errorf("%w: FSA must be 3 chars, got %q", domain.ErrInvalidInput, fsa)
    }
    if !fsaRegex.MatchString(strings.ToUpper(fsa)) {
        return fmt.Errorf("%w: FSA %q does not match %s", domain.ErrInvalidInput, fsa, fsaRegex)
    }
    return nil
}

//Generates the URL of the product information API
func (f *Fetcher) productURL(sku string) (string, error) {
    if _, err := strconv.Atoi(sku); err != nil || sku == "" {
        return "", fmt.Errorf("%w: sku %q", domain.ErrInvalidInput, sku)
    }
    u := &url.URL{
        Scheme: f.cfg.Scheme,
        Host:   f.cfg.Host,
        Path:   fmt.Sprintf(f.cfg.ProductPath, sku),
    }
    return u.String(), nil
}
//Generates the URL of the location information for a given FSA/Postal Code
func (f *Fetcher) locationsURL(fsa string) (string, error) {
    if err := validateFSA(fsa); err != nil {
        return "", err
    }
    u := &url.URL{Scheme: f.cfg.Scheme, Host: f.cfg.Host, Path: f.cfg.LocationsPath}
    q := u.Query()
    q.Set("lang", f.cfg.Language)
    q.Set("postalCode", fsa)
    u.RawQuery = q.Encode()
    return u.String(), nil
}
//Generates the URL of the availability details for each stores based on the FSA/Postal Code
func (f *Fetcher) availabilityURL(skus []string, fsa string, storeIDs []string) (string, error) {
    if len(skus) == 0 {
        return "", fmt.Errorf("%w: skus empty", domain.ErrInvalidInput)
    }
    if err := validateFSA(fsa); err != nil {
        return "", err
    }
    if len(storeIDs) == 0 {
        return "", fmt.Errorf("%w: storeIDs empty", domain.ErrInvalidInput)
    }

    u := &url.URL{Scheme: f.cfg.Scheme, Host: f.cfg.Host, Path: f.cfg.AvailabilityPath}
    q := u.Query()
    q.Set("accept", "application/vnd.bestbuy.standardproduct.v1+json")
    q.Set("accept-language", f.cfg.Language)
    q.Set("locations", strings.Join(storeIDs, "|"))
    q.Set("postalCode", fsa)
    q.Set("skus", strings.Join(skus, "|"))
    for k, v := range f.cfg.AvailabilityQ {
        if v != "" {
            q.Set(k, v)
        }
    }
    u.RawQuery = q.Encode()
    return u.String(), nil
}
