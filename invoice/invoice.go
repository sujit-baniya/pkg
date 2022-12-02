package invoice

import (
	"io"
	"strings"
	"time"
	"verify/utils"

	"github.com/jinzhu/now"
	"verify/utils/gofpdf"
)

type invoice struct {
	pdf             *gofpdf.Fpdf
	customerInvoice *CustomerInvoice
	business        *Business
	colors          *BillColor
}

type Invoice struct {
	business *Business
	colors   *BillColor
}

func New(business *Business) *Invoice {
	return &Invoice{
		business: business,
		colors: &BillColor{
			ColorLight: Color{},
			ColorDark:  Color{},
		},
	}
}

// lightText sets the font color to the light branding color from
// the config file.
func (b *invoice) lightText() {
	b.pdf.SetTextColor(
		b.colors.ColorLight.R,
		b.colors.ColorLight.G,
		b.colors.ColorLight.B,
	)
}

// darkText sets the font color to the dark branding color from
// the config file.
func (b *invoice) darkText() {
	b.pdf.SetTextColor(
		b.colors.ColorDark.R,
		b.colors.ColorDark.G,
		b.colors.ColorDark.B,
	)
}

// blackText sets the text color to black
func (b *invoice) blackText() {
	b.pdf.SetTextColor(0, 0, 0)
}

// whiteText sets the text color to black
func (b *invoice) whiteText() {
	b.pdf.SetTextColor(255, 255, 255)
}

func (b *invoice) darkDrawColor() {
	b.pdf.SetDrawColor(
		b.colors.ColorDark.R,
		b.colors.ColorDark.G,
		b.colors.ColorDark.B,
	)
}

func (b *invoice) lightFillColor() {
	b.pdf.SetFillColor(
		b.colors.ColorLight.R,
		b.colors.ColorLight.G,
		b.colors.ColorLight.B,
	)
}

func (b *invoice) fillColor(hex string) {
	if strings.Contains(hex, "#") {
		hex = hex[1:]
	}
	rgb, err := utils.Hex2RGB(hex)
	if err != nil {
		b.pdf.SetFillColor(0, 0, 0)
		return
	}
	b.pdf.SetFillColor(rgb.Red, rgb.Green, rgb.Blue)
}

// makeHeader returns the function that will be called to build
// the page header. It allows wrapping up the Fpdf instance in
// the closure.
func (b *invoice) makeHeader() func() {
	return func() {
		b.pdf.SetFont(b.business.Detail.SansFont, "", 28)
		if b.business.Detail.ImageFile != "" {
			b.pdf.ImageOptions(b.business.Detail.ImageFile, 10, 10, 70, 0, false, gofpdf.ImageOptions{}, 0, "")
		}

		// invoice Text
		b.pdf.SetXY(140, 18)
		b.text(10, 0, "Invoice", "10b981")

		// invoice Text
		b.pdf.SetXY(142, 23)
		b.pdf.SetFont(b.business.Detail.SansFont, "", 12)
		width := 20.0
		if strings.ToLower(b.customerInvoice.Detail.Status) == "unpaid" {
			b.fillColor("A8424B")
		} else if strings.ToLower(b.customerInvoice.Detail.Status) == "partially paid" {
			b.fillColor("A8863E")
			width = 38
		} else if strings.ToLower(b.customerInvoice.Detail.Status) == "paid" {
			b.fillColor("33839E")
			width = 15
		}

		b.textColor("ffffff")
		b.textFormat(width, 5, strings.ToUpper(b.customerInvoice.Detail.Status), "0", 0, "C", true, 0, "")

		// Date and invoice #
		b.pdf.SetXY(140, 35)
		b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
		b.text(20, 0, b.business.Detail.Email)

		// Date and invoice #
		b.pdf.SetXY(140, 40)
		b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
		b.text(20, 0, b.business.Detail.Person)

		// Biller Name, Address
		b.pdf.SetXY(8, 35)
		b.pdf.SetFont(b.business.Detail.SerifFont, "B", 14)
		b.text(40, 0, b.business.Detail.Name)

		b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
		b.pdf.SetXY(8, 40)
		b.text(40, 0, b.business.Detail.Address)

		// Line Break
		b.pdf.Ln(10)
		b.darkDrawColor()
		b.pdf.Line(8, 45, 200, 45)
	}
}

// makeFooter returns the function that will be called to build
// the page footer. It allows wrapping up the Fpdf instance in
// the closure.
func (b *invoice) makeFooter() func() {
	return func() {
		b.pdf.Ln(10)
		b.darkDrawColor()
		b.pdf.Line(8, 275, 200, 275)
		if b.business.Detail.FooterImageFile != "" {
			b.pdf.ImageOptions(b.business.Detail.FooterImageFile, 8.0, 277, 5, 0, false, gofpdf.ImageOptions{}, 0, "")
		}
		b.pdf.SetXY(15, 280)
		b.darkText()
		b.text(132, 0, b.business.Detail.Name)
		b.lightText()
		b.text(40, 0, "Generated: "+b.customerInvoice.Detail.GeneratedDate)
	}
}

func (b *Invoice) Create(customerInvoice *CustomerInvoice) *invoice {
	inv := &invoice{
		pdf:             gofpdf.New("P", "mm", "A4", ""),
		colors:          b.colors,
		business:        b.business,
		customerInvoice: customerInvoice,
	}

	if len(inv.customerInvoice.Detail.Date) < 1 {
		inv.customerInvoice.Detail.Date = time.Now().String()
	}

	inv.pdf.SetHeaderFunc(inv.makeHeader())
	inv.pdf.SetFooterFunc(inv.makeFooter())
	inv.pdf.AddPage()
	inv.render()
	return inv
}

func (b *invoice) RenderToFile(outFileName string) error {
	return b.pdf.OutputFileAndClose(outFileName)
}

func (b *invoice) Render(writer io.Writer) error {
	return b.pdf.Output(writer)
}

func (b *invoice) render() {
	headers := []string{"Qty", "Description", "Unit Price", "Line Total"}
	widths := []float64{16, 125.5, 25, 25}
	totalAmount := b.customerInvoice.Detail.TotalAmount
	if totalAmount == "" {
		totalAmount = b.getTotal(b.customerInvoice.Items, b.business.Tax)
	}
	b.drawBillTo(totalAmount)
	b.drawBillablesTable(headers, b.customerInvoice.Items, b.business.Tax, widths)
	if len(b.customerInvoice.Transactions) > 0 && strings.ToLower(b.customerInvoice.Detail.Status) == "paid" {
		widths = []float64{30, 35, 35, 60, 25}
		headers = []string{"Date", "Payment Method", "Transaction ID", "Description", "Amount"}
		b.drawTransactionsTable(headers, widths)
	} else {
		b.drawBankDetails()
	}
	b.notes()
}

// drawBillTo renders the invoice To part of the bill.
func (b *invoice) drawBillTo(totalAmount string) {
	// It's safe to MustParse here because we validated CLI args
	billTime := now.New(now.MustParse(b.customerInvoice.Detail.Date))

	// Default to rounding to end of month. Use exact date if specified
	invoiceDate := billTime.EndOfMonth()
	if b.customerInvoice.Detail.UseExactDate {
		invoiceDate = billTime.Time
	}
	b.pdf.Ln(5)
	b.text(0, 0, "Billed To", "10b981")
	b.pdf.SetX(5)
	b.pdf.Ln(5)
	b.text(0, 0, b.customerInvoice.Customer.Email)
	b.pdf.SetX(5)
	b.pdf.Ln(5)
	b.text(0, 0, b.customerInvoice.Customer.Name)
	b.pdf.SetX(5)
	b.pdf.Ln(5)
	b.text(0, 0, b.customerInvoice.Customer.Street)
	b.pdf.SetX(5)
	b.pdf.Ln(5)
	b.text(0, 0, b.customerInvoice.Customer.CityStateZip)
	b.pdf.SetX(5)
	b.pdf.Ln(5)
	b.text(0, 0, b.customerInvoice.Customer.Country)

	b.pdf.SetXY(80, 55)
	b.text(0, 0, "invoice Details", "10b981")
	b.pdf.SetXY(80, 60)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 12)
	b.text(20, 0, "invoice #:")
	b.text(42, 0, b.customerInvoice.Detail.InvoiceNumber)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
	b.pdf.SetXY(80, 65)
	b.text(20, 0, "Date:")
	b.text(42, 0, invoiceDate.Format("2006-01-02"))
	b.pdf.SetXY(80, 70)
	b.text(20, 0, "Due Date:")
	b.text(42, 0, invoiceDate.Format("2006-01-02"))

	b.pdf.SetXY(150, 55)
	b.textColor("10b981")
	b.textFormat(0, 0, "invoice Total", "0", 0, "R", false, 0, "")
	b.pdf.SetXY(140, 63)
	b.textColor("000000")
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 26)
	b.textFormat(0, 0, b.customerInvoice.Detail.Currency+" "+totalAmount, "0", 0, "R", false, 0, "")
}

// drawBillTable renders the summary table for the bill showing the
// department, currency, and terms.
func (b *invoice) drawBillTable(headers []string, values []string) {
	b.pdf.SetFillColor(255, 0, 0)
	b.whiteText()
	b.pdf.SetDrawColor(64, 64, 64)
	b.lightFillColor()
	b.pdf.SetLineWidth(0.3)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 10)

	baseY := b.pdf.GetY() + 20
	b.pdf.SetY(baseY)
	for _, header := range headers {
		width := float64(len(header)) * 4.9
		b.textFormat(width, 5, header, "1", 0, "C", true, 0, "")
	}

	b.pdf.Ln(5)
	b.pdf.SetFillColor(255, 255, 255)
	b.blackText()
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 8)
	for i, val := range values {
		width := float64(len(headers[i])) * 4.9
		b.textFormat(width, 4, val, "1", 0, "L", true, 0, "")
	}

}

// drawTransactionsTable renders the summary table for the bill showing the
// department, currency, and terms.
func (b *invoice) drawTransactionsTable(headers []string, widths []float64) {

	b.pdf.Ln(15)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 12)
	b.darkText()
	b.text(40, 0, "Transaction Details")

	b.pdf.SetFillColor(255, 0, 0)
	b.whiteText()
	b.pdf.SetDrawColor(64, 64, 64)
	b.lightFillColor()
	b.pdf.SetLineWidth(0.3)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 10)

	baseY := b.pdf.GetY() + 5
	b.pdf.SetY(baseY)
	for i, header := range headers {
		b.fillColor("10b981")
		b.textFormat(widths[i], 5, header, "0", 0, "R", true, 0, "")
	}

	b.pdf.Ln(5)
	b.pdf.SetFillColor(255, 255, 255)
	b.blackText()
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 8)
	for _, transaction := range b.customerInvoice.Transactions {
		for i, val := range transaction.Strings() {
			b.textFormat(widths[i], 4, val, "0", 0, "R", true, 0, "")
		}
	}

}

// drawBlanks is used to fill in the blank spaces in the table
// that precede, for example, the sub-total, tax, and total entries.
func (b *invoice) drawBlanks(billables []Item, widths []float64) {
	emptyFields := len(billables[0].Strings()) - 2
	for i := 0; i < emptyFields; i++ {
		b.textFormat(widths[i], 4, "", "", 0, "C", true, 0, "")
	}
}

// drawBillableaTable renders the table containing one line each
// for the billable items described in the YAML file.
func (b *invoice) drawBillablesTable(headers []string, billables []Item, taxDetails *TaxDetails, widths []float64) {
	b.pdf.SetFillColor(255, 0, 0)
	b.whiteText()
	b.pdf.SetDrawColor(64, 64, 64)
	b.lightFillColor()
	b.pdf.SetLineWidth(0.1)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 10)

	baseY := b.pdf.GetY() + 20
	b.pdf.SetY(baseY)
	for i, header := range headers {
		b.fillColor("10b981")
		b.textFormat(widths[i], 5, header, "0", 0, "R", true, 0, "")
	}

	b.pdf.Ln(5)
	b.pdf.SetFillColor(255, 255, 255)
	b.blackText()
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)

	// Keep the sub-total as we run through it
	var subTotal float64

	// Draw the billable items
	for _, billable := range billables {
		b.pdf.Ln(1)
		for i, val := range billable.Strings() {
			b.textFormat(widths[i], 4, val, "0", 0, "R", true, 0, "")
		}
		subTotal += billable.Total()
		b.pdf.Ln(4)
	}

	// Calculate tax
	var tax float64
	if taxDetails != nil {
		tax = subTotal * taxDetails.DefaultPercentage
	}
	total := subTotal + tax

	// Draw the Sub-Total
	b.pdf.SetDrawColor(255, 255, 255)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
	b.pdf.Ln(2)
	b.drawBlanks(billables, widths)
	subTotalText := billables[0].Currency + " " + FloatStr(subTotal)
	b.textFormat(widths[len(widths)-2], 4, "Subtotal", "1", 0, "R", true, 0, "")
	b.textFormat(widths[len(widths)-1], 4, subTotalText, "1", 0, "R", true, 0, "")

	// Handle configurable tax name
	taxName := "Tax"
	if taxDetails != nil && taxDetails.TaxName != "" {
		taxName = taxDetails.TaxName
	}

	// Draw Tax
	b.pdf.Ln(4)
	b.drawBlanks(billables, widths)
	taxText := billables[0].Currency + " " + FloatStr(tax)
	b.textFormat(widths[len(widths)-2], 4, taxName, "1", 0, "R", true, 0, "")
	b.textFormat(widths[len(widths)-1], 4, taxText, "1", 0, "R", true, 0, "")

	// Draw Total
	b.pdf.Ln(4)
	b.drawBlanks(billables, widths)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 10)
	y := b.pdf.GetY()
	x := b.pdf.GetX()
	totalText := billables[0].Currency + " " + FloatStr(total)
	b.textColor("10b981")
	b.textFormat(widths[len(widths)-2], 6, "Total", "1", 0, "R", true, 0, "")
	b.textFormat(widths[len(widths)-1], 6, totalText, "1", 0, "R", true, 0, "")
	x2 := b.pdf.GetX()

	b.pdf.SetDrawColor(64, 64, 64)
	b.pdf.Line(x, y, x2, y)
}

func (b *invoice) getTotal(billables []Item, taxDetails *TaxDetails) string {
	var subTotal float64
	for _, billable := range billables {
		subTotal += billable.Total()
	}
	var tax float64
	if taxDetails != nil {
		tax = subTotal * taxDetails.DefaultPercentage
	}
	total := subTotal + tax
	return billables[0].Currency + " " + FloatStr(total)
}

// drawBankDetails renders the table that contains the bank details.
func (b *invoice) drawBankDetails() {
	b.pdf.Ln(15)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 12)
	b.darkText()
	b.text(40, 0, "Payment Details")
	b.pdf.Ln(5)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
	b.text(40, 0, "PayPal: "+b.business.PayPal.Account, "33839E")
	b.pdf.Ln(5)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
	b.text(40, 0, "Bank Information", "33839E")
	b.pdf.Ln(5)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 8)
	headers := []string{
		"Pay By", "Bank Name", "Address", "Account Type", "Account Number",
		"Account Name", "Routing Number",
		"IBAN Code", "Sort Code (international)", "SWIFT/BIC (international)",
	}

	b.pdf.SetDrawColor(64, 64, 64)
	for _, bank := range b.business.Banks {
		for i, v := range bank.Strings() {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			b.fillColor("ffffff")
			b.textColor("10b981")
			b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
			b.textFormat(50, 5, headers[i], "0", 0, "L", true, 0, "")
			b.lightFillColor()
			b.blackText()
			b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
			b.textFormat(100, 5, v, "0", 0, "L", false, 0, "")
			b.pdf.Ln(5)
		}
		b.pdf.Ln(5)
	}
}

func (b *invoice) notes() {
	b.pdf.Ln(15)
	b.pdf.SetFont(b.business.Detail.SerifFont, "B", 10)
	b.textColor("47829e")
	b.text(40, 0, "Notes")
	b.pdf.Ln(5)
	b.pdf.SetFont(b.business.Detail.SerifFont, "", 10)
	b.text(40, 0, "Thank you for your business", "000000")
}

func (b *invoice) textColor(hex string) {
	if strings.Contains(hex, "#") {
		hex = hex[1:]
	}
	rgb, err := utils.Hex2RGB(hex)
	if err != nil {
		b.pdf.SetTextColor(0, 0, 0)
		return
	}
	b.pdf.SetTextColor(rgb.Red, rgb.Green, rgb.Blue)
}

func (b *invoice) text(x, y float64, txtStr string, hexColor ...string) {
	if len(hexColor) > 0 {
		b.textColor(hexColor[0])
	} else {
		b.textColor("000000")
	}
	unicodeToPDF := b.pdf.UnicodeTranslatorFromDescriptor("")
	b.pdf.Cell(x, y, unicodeToPDF(txtStr))
}

func (b *invoice) textFormat(w, h float64, txtStr string, borderStr string, ln int,
	alignStr string, fill bool, link int, linkStr string) {
	unicodeToPDF := b.pdf.UnicodeTranslatorFromDescriptor("")
	b.pdf.CellFormat(w, h, unicodeToPDF(txtStr), borderStr, ln, alignStr, fill, link, linkStr)
}
