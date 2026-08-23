package atoship

import (
	"encoding/json"
	"errors"
)

// Address is a postal address.
//
// State is the two-letter code for US and CA ("CA", not "California"). Country
// defaults to US when empty. Residential affects the price on most carriers, and
// getting it wrong is the usual reason a label costs more than its quote — when
// you know, say so.
type Address struct {
	Name        string `json:"name,omitempty"`
	Company     string `json:"company,omitempty"`
	Street1     string `json:"street1"`
	Street2     string `json:"street2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Country     string `json:"country,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Residential *bool  `json:"residential,omitempty"`

	// Set on addresses returned by the API.
	ID        string `json:"id,omitempty"`
	Object    string `json:"object,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// AddressRef is an address given either by ID or inline. The API accepts both in
// the same field, so this marshals to a bare string or a nested object.
//
//	atoship.AddressRef{ID: "addr_abc123"}
//	atoship.AddressRef{Inline: &atoship.Address{Street1: …}}
type AddressRef struct {
	ID     string
	Inline *Address
}

// AddressByID is shorthand for referring to a saved address.
func AddressByID(id string) AddressRef { return AddressRef{ID: id} }

// InlineAddress is shorthand for supplying an address in the request.
func InlineAddress(a Address) AddressRef { return AddressRef{Inline: &a} }

func (r AddressRef) MarshalJSON() ([]byte, error) {
	if r.Inline != nil {
		return json.Marshal(r.Inline)
	}
	if r.ID == "" {
		return nil, errors.New("atoship: AddressRef needs either ID or Inline")
	}
	return json.Marshal(r.ID)
}

func (r *AddressRef) UnmarshalJSON(b []byte) error {
	var id string
	if json.Unmarshal(b, &id) == nil {
		r.ID = id
		return nil
	}
	var a Address
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	r.Inline = &a
	return nil
}

// Parcel is what is being shipped.
//
// Dimensions matter even when they seem not to: carriers bill the greater of
// actual and dimensional weight, so omitting them on a light, bulky box quotes a
// price the carrier will not honour and the difference is billed back later.
type Parcel struct {
	Length            float64 `json:"length,omitempty"`
	Width             float64 `json:"width,omitempty"`
	Height            float64 `json:"height,omitempty"`
	Weight            float64 `json:"weight"`
	WeightUnit        string  `json:"weight_unit,omitempty"`    // oz (default), lb, g, kg
	DimensionUnit     string  `json:"dimension_unit,omitempty"` // in (default), cm
	PredefinedPackage string  `json:"predefined_package,omitempty"`
}

// ParcelRef is a parcel given either by ID or inline.
type ParcelRef struct {
	ID     string
	Inline *Parcel
}

// ParcelByID is shorthand for referring to a saved parcel.
func ParcelByID(id string) ParcelRef { return ParcelRef{ID: id} }

// InlineParcel is shorthand for supplying parcel dimensions in the request.
func InlineParcel(p Parcel) ParcelRef { return ParcelRef{Inline: &p} }

func (r ParcelRef) MarshalJSON() ([]byte, error) {
	if r.Inline != nil {
		return json.Marshal(r.Inline)
	}
	if r.ID == "" {
		return nil, errors.New("atoship: ParcelRef needs either ID or Inline")
	}
	return json.Marshal(r.ID)
}

func (r *ParcelRef) UnmarshalJSON(b []byte) error {
	var id string
	if json.Unmarshal(b, &id) == nil {
		r.ID = id
		return nil
	}
	var p Parcel
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	r.Inline = &p
	return nil
}

// ShipmentOptions are extras requested at quote time. Whatever you ask for here
// must also be sent when buying, or the label is bought without it at a price
// quoted with it.
type ShipmentOptions struct {
	InsuranceAmount       float64            `json:"insurance_amount,omitempty"`
	SignatureConfirmation bool               `json:"signature_confirmation,omitempty"`
	SaturdayDelivery      bool               `json:"saturday_delivery,omitempty"`
	LithiumBatteries      *LithiumBatteries  `json:"lithium_batteries,omitempty"`
	DryIce                *DryIce            `json:"dry_ice,omitempty"`
	Alcohol               *Alcohol           `json:"alcohol,omitempty"`
	Incoterm              string             `json:"incoterm,omitempty"`
	LabelCustomFields     *LabelCustomFields `json:"label_custom_fields,omitempty"`
}

// LithiumBatteries declares a hazmat battery shipment.
type LithiumBatteries struct {
	Type        string `json:"type,omitempty"`
	PackingType string `json:"packing_type,omitempty"`
}

// DryIce declares dry ice, which is a regulated material.
type DryIce struct {
	Weight     float64 `json:"weight,omitempty"`
	WeightUnit string  `json:"weight_unit,omitempty"`
}

// Alcohol declares an alcohol shipment and the licence it ships under.
type Alcohol struct {
	RecipientType string `json:"recipient_type,omitempty"`
	LicenseNumber string `json:"license_number,omitempty"`
	LicenseState  string `json:"license_state,omitempty"`
}

// LabelCustomFields are printed on the label where the carrier supports it.
// Where it does not, the value is stored with the label and returned by the API
// but does not appear on the label itself.
type LabelCustomFields struct {
	SKU           string `json:"sku,omitempty"`
	UPC           string `json:"upc,omitempty"`
	Reference1    string `json:"reference1,omitempty"`
	Reference2    string `json:"reference2,omitempty"`
	Reference3    string `json:"reference3,omitempty"`
	InvoiceNumber string `json:"invoice_number,omitempty"`
	PrintBarcode  bool   `json:"print_barcode,omitempty"`
	BarcodeType   string `json:"barcode_type,omitempty"`
}

// Rate is one priced service for a shipment. Rate.ID is what Labels.Create
// takes; it is valid for a limited time, so quote and buy in the same session
// rather than storing rates.
type Rate struct {
	ID                     string       `json:"id"`
	Object                 string       `json:"object"`
	Carrier                string       `json:"carrier"`
	Service                string       `json:"service"`
	ServiceCode            string       `json:"service_code"`
	Rate                   float64      `json:"rate"`
	Currency               string       `json:"currency"`
	Zone                   string       `json:"zone"`
	DeliveryDays           int          `json:"delivery_days"`
	DeliveryDate           string       `json:"delivery_date"`
	DeliveryDateGuaranteed bool         `json:"delivery_date_guaranteed"`
	ListRate               float64      `json:"list_rate"`
	RetailRate             float64      `json:"retail_rate"`
	CarrierAccountID       string       `json:"carrier_account_id"`
	Includes               RateIncludes `json:"includes"`
	CreatedAt              string       `json:"created_at"`
}

// RateIncludes says which requested options this price already covers.
type RateIncludes struct {
	Insurance bool `json:"insurance"`
	Signature bool `json:"signature"`
	Saturday  bool `json:"saturday"`
}

// Label is a shipping label.
type Label struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Status           string `json:"status"`
	TrackingNumber   string `json:"tracking_number"`
	TrackingPending  bool   `json:"tracking_pending"`
	Carrier          string `json:"carrier"`
	CarrierAccountID string `json:"carrier_account_id"`
	Service          string `json:"service"`
	ServiceCode      string `json:"service_code"`
	// Rate is what you were charged. Nil on a small number of labels created
	// before the price was recorded — unknown, not free.
	Rate        *float64 `json:"rate"`
	Currency    string   `json:"currency"`
	LabelURL    string   `json:"label_url"`
	LabelFormat string   `json:"label_format"`
	FromAddress Address  `json:"from_address"`
	ToAddress   Address  `json:"to_address"`
	Parcel      Parcel   `json:"parcel"`
	Reference   string   `json:"reference"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	PurchasedAt string   `json:"purchased_at"`
}

// listEnvelope is the shape of every paginated response.
type listEnvelope[T any] struct {
	Object     string     `json:"object"`
	Mode       string     `json:"mode"`
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
