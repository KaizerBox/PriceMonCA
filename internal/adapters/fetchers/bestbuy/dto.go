package bestbuy

type productDTO struct {
    SKU          string  `json:"sku"`
    Name         string  `json:"name"`
    RegularPrice float64 `json:"regularPrice"`
    SalePrice    float64 `json:"salePrice"`
    SellerID     string  `json:"sellerId"`
    Currency     string  `json:"currencyCode"`
    ProductURL   string  `json:"productUrl"`
}

type storeDTO struct {
    StoreID      string `json:"locationID"`
    StoreName    string `json:"name"`
    StoreAddress string `json:"address1"`
}

type locationsDTO struct {
    Locations []storeDTO `json:"locations"`
}

func (l locationsDTO) StoreIDs() []string {
    ids := make([]string, 0, len(l.Locations))
    for _, s := range l.Locations {
        ids = append(ids, s.StoreID)
    }
    return ids
}

func (l locationsDTO) IndexByID() map[string]storeDTO {
    m := make(map[string]storeDTO, len(l.Locations))
    for _, s := range l.Locations {
        m[s.StoreID] = s
    }
    return m
}

type pickupDetail struct {
    StoreID             string `json:"locationKey"`
    QuantityOnHand      int    `json:"quantityOnHand"`
    HasInventory        bool   `json:"hasInventory"`
    SupportsFulfillment bool   `json:"supportsFulfillment"`
    IsReservable        bool   `json:"isReservable"`
}

type pickupAvailability struct {
    Status  string         `json:"status"`
    Details []pickupDetail `json:"locations"`
}

type shippingDetail struct {
    CarrierName           string `json:"carrierName"`
    DeliveryDate          string `json:"deliveryDate"`
    DeliveryDateExpiresOn string `json:"deliveryDateExpiresOn"`
    DeliveryDatePrecision string `json:"deliveryDatePrecision"`
}

type shippingAvailability struct {
    Status            string           `json:"status"`
    QuantityRemaining int              `json:"quantityRemaining"`
    IsBackorderable   bool             `json:"isBackorderable"`
    Details           []shippingDetail `json:"levelsOfServices"`
}

type productAvailabilityDTO struct {
    SKU      string               `json:"sku"`
    SellerID string               `json:"sellerId"`
    Pickup   pickupAvailability   `json:"pickup"`
    Shipping shippingAvailability `json:"shipping"`
}

type availabilityDTO struct {
    Availabilities []productAvailabilityDTO `json:"availabilities"`
}

func (a availabilityDTO) BySKU() map[string]productAvailabilityDTO {
    m := make(map[string]productAvailabilityDTO, len(a.Availabilities))
    for _, p := range a.Availabilities {
        m[p.SKU] = p
    }
    return m
}
