package invoice

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"regexp"
	"strconv"
)

type BusinessDetails struct {
	Name            string `yaml:"name" json:"name"`
	Person          string `yaml:"person" json:"person"`
	Email           string `yaml:"email" json:"email"`
	Address         string `yaml:"address" json:"address"`
	ImageFile       string `yaml:"image_file" json:"image_file"`
	FooterImageFile string `yaml:"footer_image_file" json:"footer_image_file"`
	SansFont        string `yaml:"sans_font" json:"sans_font"`
	SerifFont       string `yaml:"serif_font" json:"serif_font"`
}

type Detail struct {
	Department    string `yaml:"department" json:"department"`
	InvoiceNumber string `yaml:"invoice_number" json:"invoice_number"`
	Currency      string `yaml:"currency" json:"currency"`
	PaymentTerms  string `yaml:"payment_terms" json:"payment_terms"`
	DueDate       string `yaml:"due_date" json:"due_date"`
	Date          string `yaml:"date" json:"date"`
	Status        string `yaml:"status" json:"status"`
	UseExactDate  bool   `yaml:"use_exact_date" json:"use_exact_date"`
	GeneratedDate string `yaml:"generated_date" json:"generated_date"`
	TotalAmount   string `yaml:"total_amount" json:"total_amount"`
}

func (b *Detail) Strings() []string {
	return []string{
		b.Department, b.InvoiceNumber, b.Currency, b.PaymentTerms, b.DueDate,
	}
}

type Transaction struct {
	Date          string  `json:"date" yaml:"date"`
	PaymentMethod string  `yaml:"payment_method" json:"payment_method"`
	ID            string  `json:"id" yaml:"id"`
	Description   string  `yaml:"description" json:"description"`
	Amount        float64 `yaml:"amount" json:"amount"`
	Currency      string  `yaml:"currency" json:"currency"`
}

func (b *Transaction) Strings() []string {
	return []string{
		b.Date, b.PaymentMethod, b.ID, b.Description, b.Currency + FloatStr(b.Amount),
	}
}

type Customer struct {
	Email        string `yaml:"email" json:"email"`
	Name         string `yaml:"name" json:"name"`
	Street       string `yaml:"street" json:"street"`
	CityStateZip string `yaml:"city_state_zip" json:"city_state_zip"`
	Country      string `yaml:"country" json:"country"`
}

type Item struct {
	Quantity    float64 `yaml:"quantity" json:"quantity"`
	Description string  `yaml:"description" json:"description"`
	UnitPrice   float64 `yaml:"unit_price" json:"unit_price"`
	Currency    string  `yaml:"currency" json:"currency"`
}

func (b *Item) Total() float64 {
	return b.UnitPrice * b.Quantity
}

func (b *Item) Strings() []string {
	return []string{
		strconv.FormatFloat(b.Quantity, 'f', 2, 64),
		b.Description,
		b.Currency + " " + FloatStr(b.UnitPrice),
		b.Currency + " " + FloatStr(b.Total()),
	}
}

type TaxDetails struct {
	DefaultPercentage float64 `yaml:"default_percentage" json:"default_percentage"`
	TaxName           string  `yaml:"tax_name" json:"tax_name"`
}

type BankDetails struct {
	TransferType  string `yaml:"transfer_type" json:"transfer_type"`
	Name          string `yaml:"name" json:"name"`
	Address       string `yaml:"address" json:"address"`
	AccountType   string `yaml:"account_type" json:"account_type"`
	AccountNumber string `yaml:"account_number" json:"account_number"`
	AccountName   string `yaml:"account_name" json:"account_name"`
	RoutingNumber string `yaml:"routing_number" json:"routing_number"`
	IBAN          string `yaml:"iban" json:"iban"`
	SortCode      string `yaml:"sort_code" json:"sort_code"`
	SWIFTBIC      string `yaml:"swift_bic" json:"swift_bic"`
}

func (b *BankDetails) Strings() []string {
	return []string{
		b.TransferType, b.Name, b.Address, b.AccountType, b.AccountNumber, b.AccountName, b.RoutingNumber, b.IBAN, b.SortCode, b.SWIFTBIC,
	}
}

type Color struct {
	R int
	G int
	B int
}

type BillColor struct {
	ColorLight Color `yaml:"color_light" json:"color_light"`
	ColorDark  Color `yaml:"color_dark" json:"color_dark"`
}

type Paypal struct {
	Account string `yaml:"account" json:"account"`
}

type Business struct {
	Detail *BusinessDetails `yaml:"business" json:"business"`
	Tax    *TaxDetails      `yaml:"tax" json:"tax"`
	Banks  []BankDetails    `yaml:"banks" json:"banks"`
	PayPal *Paypal          `yaml:"paypal" json:"paypal"`
}

type CustomerInvoice struct {
	Detail       *Detail       `yaml:"detail" json:"detail"`
	Customer     *Customer     `yaml:"customer" json:"customer"`
	Items        []Item        `yaml:"items" json:"items"`
	Transactions []Transaction `yaml:"transactions" json:"transactions"`
}

// FloatStr takes a float and gives back a monetary, human-formatted
// value.
var r = regexp.MustCompile("-?[0-9,]+.[0-9]{2}")

func FloatStr(f float64) string {
	roundedFloat := math.Round(f*100) / 100
	p := message.NewPrinter(language.English)
	results := r.FindAllString(p.Sprintf("%f", roundedFloat), 1)

	if len(results) < 1 {
		panic("got some ridiculous number that has no decimals")
	}

	return results[0]
}
