package money

import (
	"fmt"
	"github.com/sujit-baniya/framework/facades"
	"net/http"
	"time"

	"encoding/json"
)

var CurrencySymbol = map[string]string{
	"AFN": "؋",
	"ALL": "L",
	"DZD": "دج",
	"AOA": "Kz",
	"ARS": "$",
	"AMD": "֏",
	"AWG": "ƒ",
	"AZN": "₼",
	"BHD": ".د.ب",
	"BSD": "$",
	"BDT": "৳",
	"BBD": "$",
	"BYR": "Br",
	"BZD": "BZ$",
	"BMD": "$",
	"BTN": "Nu.",
	"BOB": "$b",
	"BAM": "KM",
	"BWP": "P",
	"BRL": "R$",
	"BND": "$",
	"BGN": "лв",
	"MMK": "K",
	"BIF": "FBu",
	"KHR": "៛",
	"CAD": "$",
	"CVE": "$",
	"KYD": "$",
	"CLP": "$",
	"CNY": "¥",
	"COP": "$",
	"KMF": "CF",
	"CDF": "FC",
	"CRC": "₡",
	"HRK": "kn",
	"CUP": "₱",
	"CZK": "Kč",
	"DJF": "Fdj",
	"DOP": "RD$",
	"EGP": "£",
	"ERN": "Nfk",
	"ETB": "Br",
	"FKP": "£",
	"XAF": "FCFA",
	"GMD": "D",
	"GEL": "₾",
	"GHS": "GH₵",
	"DKK": "kr",
	"GTQ": "Q",
	"GNF": "FG",
	"GYD": "$",
	"HTG": "G",
	"HNL": "L",
	"HKD": "$",
	"HUF": "Ft",
	"ISK": "kr",
	"INR": "₹",
	"IDR": "Rp",
	"IRR": "﷼",
	"IQD": "ع.د",
	"ILS": "₪",
	"JMD": "J$",
	"JPY": "¥",
	"JOD": "JD",
	"KZT": "лв",
	"KES": "KSh",
	"KWD": "KD",
	"KGS": "лв",
	"LAK": "₭",
	"LBP": "£",
	"LSL": "M",
	"LRD": "$",
	"LYD": "LD",
	"LTL": "Lt",
	"MOP": "MOP$",
	"MKD": "ден",
	"MGA": "Ar",
	"MWK": "MK",
	"MYR": "RM",
	"MVR": "Rf",
	"MRO": "UM",
	"MUR": "₨",
	"MXN": "$",
	"MDL": "lei",
	"MNT": "₮",
	"MZN": "MT",
	"NAD": "$",
	"NPR": "₨",
	"NIO": "C$",
	"NGN": "₦",
	"KPW": "₩",
	"OMR": "﷼",
	"PKR": "₨",
	"PAB": "B/.",
	"PGK": "K",
	"PYG": "Gs",
	"PEN": "S/.",
	"PHP": "₱",
	"PLN": "zł",
	"QAR": "﷼",
	"RON": "lei",
	"RUB": "₽",
	"RWF": "R₣",
	"SHP": "£",
	"STD": "Db",
	"XCD": "$",
	"WST": "WS$",
	"SAR": "﷼",
	"RSD": "Дин.",
	"SCR": "₨",
	"SLL": "Le",
	"SGD": "$",
	"ANG": "ƒ",
	"SBD": "$",
	"ZAR": "R",
	"GBP": "£",
	"KRW": "₩",
	"SSP": "£",
	"EUR": "€",
	"LKR": "₨",
	"SDG": "ج.س.",
	"SRD": "$",
	"NOK": "kr",
	"SZL": "E",
	"SEK": "kr",
	"CHF": "CHF",
	"SYP": "£",
	"TWD": "NT$",
	"TJS": "SM",
	"TZS": "TSh",
	"THB": "฿",
	"XOF": "CFA",
	"NZD": "$",
	"TOP": "T$",
	"TTD": "TT$",
	"TND": "د.ت",
	"TRY": "₺",
	"TMT": "T",
	"AUD": "$",
	"UGX": "USh",
	"UAH": "₴",
	"UYU": "$U",
	"UZS": "лв",
	"VUV": "VT",
	"VEF": "Bs",
	"VND": "₫",
	"USD": "$",
	"XPF": "₣",
	"MAD": "MAD",
	"YER": "﷼",
	"ZMK": "﷼",
	"ZWL": "﷼",
}

type CurrencyRate struct {
	Base   string  `json:"base"`
	Target string  `json:"target"`
	Rate   Decimal `json:"rate"`
	Date   string  `json:"date"`
}

type CurrencyResponse struct {
	Date  string             `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]Decimal `json:"rates"`
}

var ExchangeRateURL = "https://api.exchangerate.host"
var httpClient = http.Client{}

func GetCurrencySymbol(currency string) string {
	return CurrencySymbol[currency]
}

func Convert(amount Decimal, baseCurrency string, targetCurrency string) Decimal {
	return GetLatestRate(baseCurrency, targetCurrency).Convert(amount)
}

func GetLatestRate(baseCurrency string, targetCurrency string) *CurrencyRate {
	c := &CurrencyRate{
		Base:   baseCurrency,
		Target: targetCurrency,
	}
	c.getLatestCachedRate()
	if c.Rate.IsZero() {
		c.request("latest")
	}
	return c
}

func GetHistoryRate(baseCurrency string, targetCurrency string, date string) *CurrencyRate {
	c := &CurrencyRate{
		Base:   baseCurrency,
		Target: targetCurrency,
	}
	c.getHistoryCachedRate(date)
	if c.Rate.IsZero() {
		c.request(date)
	}
	return c
}

func (c *CurrencyRate) request(date string) {
	var cr CurrencyResponse
	url := ExchangeRateURL + "/" + date + "?base=" + c.Base + "&symbols=" + c.Target
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	json.NewDecoder(resp.Body).Decode(&cr)
	if val, ok := cr.Rates[c.Target]; ok {
		c.Rate = val
	}
	c.cache(date)
}

func (c *CurrencyRate) getLatestCachedRate() {
	c.getFromCache("latest")
}

func (c *CurrencyRate) Convert(amount Decimal) Decimal {
	return c.Rate.Mul(amount)
}

func (c *CurrencyRate) SetLatestCachedRate() {
	c.cache("latest")
}

func (c *CurrencyRate) getHistoryCachedRate(date string) {
	c.getFromCache(date)
}

func (c *CurrencyRate) SetHistoryCachedRate(date string) {
	c.cache(date)
}

func (c *CurrencyRate) getFromCache(date string) {
	key := c.Base + "_" + c.Target + "_" + date
	res, _ := facades.Memory.Get(key)
	if res != nil {
		json.Unmarshal(res, &c)
	}
}

func (c *CurrencyRate) cache(date string) {
	c.Date = date
	key := c.Base + "_" + c.Target + "_" + date
	val, _ := json.Marshal(&c)
	facades.Memory.Set(key, val, 24*time.Hour)
}
