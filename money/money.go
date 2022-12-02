package money

type Money struct {
	Amount   Decimal
	Currency string
}

func (m *Money) ToCurrency(currency string) Decimal {
	return Convert(m.Amount, m.Currency, currency)
}
