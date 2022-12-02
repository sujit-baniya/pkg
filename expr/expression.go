package expr

import (
	"fmt"
	"regexp"
	"strings"

	"mvdan.cc/xurls/v2"
)

// LinkRegex Regular expression patterns
var LinkRegex = xurls.Relaxed()
var (
	DatePattern              = `(?i)(?:[0-3]?\d(?:st|nd|rd|th)?\s+(?:of\s+)?(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)|(?:jan\.?|january|feb\.?|february|mar\.?|march|apr\.?|april|may|jun\.?|june|jul\.?|july|aug\.?|august|sep\.?|september|oct\.?|october|nov\.?|november|dec\.?|december)\s+[0-3]?\d(?:st|nd|rd|th)?)(?:\,)?\s*(?:\d{4})?|[0-3]?\d[-\./][0-3]?\d[-\./]\d{2,4}`
	TimePattern              = `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`
	PhonePattern             = `(?:(?:\+?\d{1,3}[-.\s*]?)?(?:\(?\d{3}\)?[-.\s*]?)?\d{3}[-.\s*]?\d{4,6})|(?:(?:(?:\(\+?\d{2}\))|(?:\+?\d{2}))\s*\d{2}\s*\d{3}\s*\d{4})`
	PhonesWithExtsPattern    = `(?i)(?:(?:\+?1\s*(?:[.-]\s*)?)?(?:\(\s*(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9])\s*\)|(?:[2-9]1[02-9]|[2-9][02-8]1|[2-9][02-8][02-9]))\s*(?:[.-]\s*)?)?(?:[2-9]1[02-9]|[2-9][02-9]1|[2-9][02-9]{2})\s*(?:[.-]\s*)?(?:[0-9]{4})(?:\s*(?:#|x\.?|ext\.?|extension)\s*(?:\d+)?)`
	LinkPattern              = LinkRegex.String()
	EmailPattern             = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	IPv4Pattern              = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	IPv6Pattern              = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	IPPattern                = IPv4Pattern + `|` + IPv6Pattern
	NotKnownPortPattern      = `6[0-5]{2}[0-3][0-5]|[1-5][\d]{4}|[2-9][\d]{3}|1[1-9][\d]{2}|10[3-9][\d]|102[4-9]`
	PricePattern             = `[$]\s?[+-]?[0-9]{1,3}(?:(?:,?[0-9]{3}))*(?:\.[0-9]{1,2})?`
	HexColorPattern          = `(?:#?([0-9a-fA-F]{6}|[0-9a-fA-F]{3}))`
	CreditCardPattern        = `(?:(?:(?:\d{4}[- ]?){3}\d{4}|\d{15,16}))`
	VISACreditCardPattern    = `4\d{3}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	MCCreditCardPattern      = `5[1-5]\d{2}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}`
	BtcAddressPattern        = `[13][a-km-zA-HJ-NP-Z1-9]{25,34}`
	StreetAddressPattern     = `\d{1,4} [\w\s]{1,20}(?:street|st|avenue|ave|road|rd|highway|hwy|square|sq|trail|trl|drive|dr|court|ct|park|parkway|pkwy|circle|cir|boulevard|blvd)\W?`
	ZipCodePattern           = `\b\d{5}(?:[-\s]\d{4})?\b`
	PoBoxPattern             = `(?i)P\.? ?O\.? Box \d+`
	SSNPattern               = `(?:\d{3}-\d{2}-\d{4})`
	MD5HexPattern            = `[0-9a-fA-F]{32}`
	SHA1HexPattern           = `[0-9a-fA-F]{40}`
	SHA256HexPattern         = `[0-9a-fA-F]{64}`
	GUIDPattern              = `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`
	ISBN13Pattern            = `(?:[\d]-?){12}[\dxX]`
	ISBN10Pattern            = `(?:[\d]-?){9}[\dxX]`
	MACAddressPattern        = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
	IBANPattern              = `[A-Z]{2}\d{2}[A-Z0-9]{4}\d{7}([A-Z\d]?){0,16}`
	GitRepoPattern           = `((git|ssh|http(s)?)|(git@[\w\.]+))(:(\/\/)?)([\w\.@\:/\-~]+)(\.git)(\/)?`
	PlaceHolderPattern       = `<(\w+([\s-_]\w+)*)(:([a-zA-Z_]+))?>`
	OutputPlaceHolderPattern = `<\w+([\s-_]\w+)*>`

	WordPattern   = `[a-zA-Z]+`
	IntPattern    = `[+\-]?[0-9]+`
	FloatPattern  = `[+\-]?\d[\d,]*[\.]?[\d{2}]*`
	StringPattern = `.+`
)

// Compiled regular expressions
var (
	DateRegex              = regexp.MustCompile(DatePattern)
	TimeRegex              = regexp.MustCompile(TimePattern)
	PhoneRegex             = regexp.MustCompile(PhonePattern)
	PhonesWithExtsRegex    = regexp.MustCompile(PhonesWithExtsPattern)
	EmailRegex             = regexp.MustCompile(EmailPattern)
	IPv4Regex              = regexp.MustCompile(IPv4Pattern)
	IPv6Regex              = regexp.MustCompile(IPv6Pattern)
	IPRegex                = regexp.MustCompile(IPPattern)
	NotKnownPortRegex      = regexp.MustCompile(NotKnownPortPattern)
	PriceRegex             = regexp.MustCompile(PricePattern)
	HexColorRegex          = regexp.MustCompile(HexColorPattern)
	CreditCardRegex        = regexp.MustCompile(CreditCardPattern)
	BtcAddressRegex        = regexp.MustCompile(BtcAddressPattern)
	StreetAddressRegex     = regexp.MustCompile(StreetAddressPattern)
	ZipCodeRegex           = regexp.MustCompile(ZipCodePattern)
	PoBoxRegex             = regexp.MustCompile(PoBoxPattern)
	SSNRegex               = regexp.MustCompile(SSNPattern)
	MD5HexRegex            = regexp.MustCompile(MD5HexPattern)
	SHA1HexRegex           = regexp.MustCompile(SHA1HexPattern)
	SHA256HexRegex         = regexp.MustCompile(SHA256HexPattern)
	GUIDRegex              = regexp.MustCompile(GUIDPattern)
	ISBN13Regex            = regexp.MustCompile(ISBN13Pattern)
	ISBN10Regex            = regexp.MustCompile(ISBN10Pattern)
	VISACreditCardRegex    = regexp.MustCompile(VISACreditCardPattern)
	MCCreditCardRegex      = regexp.MustCompile(MCCreditCardPattern)
	MACAddressRegex        = regexp.MustCompile(MACAddressPattern)
	IBANRegex              = regexp.MustCompile(IBANPattern)
	GitRepoRegex           = regexp.MustCompile(GitRepoPattern)
	PlaceHolderRegex       = regexp.MustCompile(PlaceHolderPattern)
	OutputPlaceHolderRegex = regexp.MustCompile(OutputPlaceHolderPattern)
	WordRegex              = regexp.MustCompile(WordPattern)
	IntRegex               = regexp.MustCompile(IntPattern)
	FloatRegex             = regexp.MustCompile(FloatPattern)
	StringRegex            = regexp.MustCompile(StringPattern)
)

var PatternsMap = map[string]string{
	"date":           DatePattern,
	"time":           TimePattern,
	"phone":          PhonePattern,
	"phone_ext":      PhonesWithExtsPattern,
	"link":           LinkPattern,
	"email":          EmailPattern,
	"ip4":            IPv4Pattern,
	"ip6":            IPv6Pattern,
	"ip":             IPPattern,
	"price":          PricePattern,
	"hex_color":      HexColorPattern,
	"cc":             CreditCardPattern,
	"visa_cc":        VISACreditCardPattern,
	"mc_cc":          MCCreditCardPattern,
	"btc_address":    BtcAddressPattern,
	"street_address": StreetAddressPattern,
	"zip_code":       ZipCodePattern,
	"po_box":         PoBoxPattern,
	"ssn":            SSNPattern,
	"md5":            MD5HexPattern,
	"sha1":           SHA1HexPattern,
	"sha256":         SHA256HexPattern,
	"guid":           GUIDPattern,
	"isbn_10":        ISBN13Pattern,
	"isbn_13":        ISBN10Pattern,
	"mac_address":    MACAddressPattern,
	"iban":           IBANPattern,
	"git_repo":       GitRepoPattern,
	"word":           WordPattern,
	"int":            IntPattern,
	"integer":        IntPattern,
	"float":          FloatPattern,
	"string":         StringPattern,
}

var RegexMap = map[string]*regexp.Regexp{
	"date":           DateRegex,
	"time":           TimeRegex,
	"phone":          PhoneRegex,
	"phone_ext":      PhonesWithExtsRegex,
	"link":           LinkRegex,
	"email":          EmailRegex,
	"ip4":            IPv4Regex,
	"ip6":            IPv6Regex,
	"ip":             IPRegex,
	"price":          PriceRegex,
	"cc":             CreditCardRegex,
	"visa_cc":        VISACreditCardRegex,
	"mc_cc":          MCCreditCardRegex,
	"btc_address":    BtcAddressRegex,
	"street_address": StreetAddressRegex,
	"zip_code":       ZipCodeRegex,
	"po_box":         PoBoxRegex,
	"ssn":            SSNRegex,
	"md5":            MD5HexRegex,
	"sha1":           SHA1HexRegex,
	"sha256":         SHA256HexRegex,
	"guid":           GUIDRegex,
	"isbn_10":        ISBN13Regex,
	"isbn_13":        ISBN10Regex,
	"mac_address":    MACAddressRegex,
	"iban":           IBANRegex,
	"git_repo":       GitRepoRegex,
}

// placeholderReplacer converts our input pattern string into a regular expression string.
func placeholderReplacer(inString string) string {
	// Expected input format: <name:type>
	// Output format: "(?P<name>...)" where ... is the regular expression for that type.

	matches := PlaceHolderRegex.FindStringSubmatch(inString)
	if matches == nil {
		return inString
	}
	name := matches[1]
	dataPattern := StringPattern
	if len(matches) > 3 {
		if pattern, found := PatternsMap[matches[3]]; found {
			dataPattern = pattern
		}
	}
	return fmt.Sprintf("(?P<%s>%s)", name, dataPattern)
}

// Transform pulls values out of the data string using inPattern, and then writes those values
// to a new string using the outPattern.
func Transform(inPattern string, outPattern string, data string) (string, error) {
	valueMap, err := Parse(data, inPattern)
	if err != nil {
		return "", err
	}
	return Replace(valueMap, outPattern)
}

func ParseMultiple(data string, patterns ...string) map[string][]string {
	dataList := make(map[string][]string)
	if len(patterns) == 0 {
		for pattern, rex := range RegexMap {
			tmp := match(data, rex)
			if len(tmp) > 0 {
				dataList[pattern] = tmp
			}
		}
		return dataList
	}
	for _, pattern := range patterns {
		if pattern == "email" {
			dataList[pattern] = ParseEmails(data)
		} else if pattern == "link" {
			dataList[pattern] = ParseLinks(data)
		} else if val, ok := RegexMap[pattern]; ok {
			tmp := match(data, val)
			if len(tmp) > 0 {
				dataList[pattern] = tmp
			}
		}
	}
	return dataList
}

func Parse(data, pattern string) (map[string]string, error) {
	regexpPatternInString := PlaceHolderRegex.ReplaceAllStringFunc(pattern, placeholderReplacer)
	regexpPatternIn := regexp.MustCompile(regexpPatternInString)
	// let's build a data map
	valueMap := make(map[string]string)
	for _, name := range regexpPatternIn.SubexpNames() {
		if name == "" || name == "skip" {
			continue
		}
		valueMap[name] = ""
	}
	// Get the data out of the input string
	matches := regexpPatternIn.FindStringSubmatch(data)
	if matches == nil {
		return valueMap, fmt.Errorf("data did not match input pattern")
	}

	// let's build a data map
	for _, name := range regexpPatternIn.SubexpNames() {
		if name == "" || name == "skip" {
			continue
		}
		valueMap[name] = matches[regexpPatternIn.SubexpIndex(name)]
	}
	return valueMap, nil
}

func Replace(data map[string]string, pattern string) (string, error) {
	var errorOccurred error = nil

	output := OutputPlaceHolderRegex.ReplaceAllStringFunc(pattern, replacer(data, &errorOccurred))
	if errorOccurred != nil {
		return "", errorOccurred
	}

	return output, nil
}

func replacer(data map[string]string, err *error) func(n string) string {
	return func(n string) string {
		name := n[1 : len(n)-1]
		val, found := data[name]
		if !found {
			return ""
		}
		return val
	}
}

func match(text string, regex *regexp.Regexp) []string {
	parsed := regex.FindAllString(text, -1)
	return parsed
}

// ParseDate finds all date strings
func ParseDate(text string) []string {
	return match(text, DateRegex)
}

// ParseTime finds all time strings
func ParseTime(text string) []string {
	return match(text, TimeRegex)
}

// ParsePhones finds all phone numbers
func ParsePhones(text string) []string {
	return match(text, PhoneRegex)
}

// ParsePhonesWithExts finds all phone numbers with ext
func ParsePhonesWithExts(text string) []string {
	return match(text, PhonesWithExtsRegex)
}

// ParseLinks finds all link strings
func ParseLinks(str string) []string {
	urls := LinkRegex.FindAllString(str, -1)
	return removeURLSchemeWithNoAuthority(urls).Urls
}

type LinkResponse struct {
	Urls      []string `json:"urls"`
	Emails    []string `json:"emails"`
	TotalUrls int      `json:"total_urls"`
}

func removeURLSchemeWithNoAuthority(urls []string) LinkResponse {
	var urlSlice []string
	var emailSlice []string
	totalUrls := 0
	for _, url := range urls {
		if isEmail(url) {
			emailSlice = append(emailSlice, url)
			totalUrls++
			continue
		}
		foundNoAuthorityScheme := false
		for _, scheme := range xurls.SchemesNoAuthority {
			if strings.Contains(url, scheme) {
				foundNoAuthorityScheme = true
				break
			}
		}
		if !foundNoAuthorityScheme {
			urlSlice = append(urlSlice, url)
			totalUrls++
		}
	}
	return LinkResponse{
		Urls:      urlSlice,
		Emails:    emailSlice,
		TotalUrls: totalUrls,
	}
}

func isEmail(email string) bool {
	i := strings.LastIndexByte(email, '@')
	if i != -1 && !strings.Contains(email, `://`) {
		return true
	}
	return false
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// ParseEmails finds all email strings
func ParseEmails(text string) []string {
	urls := LinkRegex.FindAllString(text, -1)
	return removeURLSchemeWithNoAuthority(urls).Emails
}

// ParseIPv4s finds all IPv4 addresses
func ParseIPv4s(text string) []string {
	return match(text, IPv4Regex)
}

// ParseIPv6s finds all IPv6 addresses
func ParseIPv6s(text string) []string {
	return match(text, IPv6Regex)
}

// ParseIPs finds all IP addresses (both IPv4 and IPv6)
func ParseIPs(text string) []string {
	return match(text, IPRegex)
}

// ParseNotKnownPorts finds all not-known port numbers
func ParseNotKnownPorts(text string) []string {
	return match(text, NotKnownPortRegex)
}

// ParsePrices finds all price strings
func ParsePrices(text string) []string {
	return match(text, PriceRegex)
}

// ParseHexColors finds all hex color values
func ParseHexColors(text string) []string {
	return match(text, HexColorRegex)
}

// ParseCreditCards finds all credit card numbers
func ParseCreditCards(text string) []string {
	return match(text, CreditCardRegex)
}

// ParseBtcAddresses finds all bitcoin addresses
func ParseBtcAddresses(text string) []string {
	return match(text, BtcAddressRegex)
}

// ParseStreetAddresses finds all street addresses
func ParseStreetAddresses(text string) []string {
	return match(text, StreetAddressRegex)
}

// ParseZipCodes finds all zip codes
func ParseZipCodes(text string) []string {
	return match(text, ZipCodeRegex)
}

// ParsePoBoxes finds all po-box strings
func ParsePoBoxes(text string) []string {
	return match(text, PoBoxRegex)
}

// ParseSSNs finds all SSN strings
func ParseSSNs(text string) []string {
	return match(text, SSNRegex)
}

// ParseMD5Hexes finds all MD5 hex strings
func ParseMD5Hexes(text string) []string {
	return match(text, MD5HexRegex)
}

// ParseSHA1Hexes finds all SHA1 hex strings
func ParseSHA1Hexes(text string) []string {
	return match(text, SHA1HexRegex)
}

// ParseSHA256Hexes finds all SHA256 hex strings
func ParseSHA256Hexes(text string) []string {
	return match(text, SHA256HexRegex)
}

// ParseGUIDs finds all GUID strings
func ParseGUIDs(text string) []string {
	return match(text, GUIDRegex)
}

// ParseISBN13s finds all ISBN13 strings
func ParseISBN13s(text string) []string {
	return match(text, ISBN13Regex)
}

// ParseISBN10s finds all ISBN10 strings
func ParseISBN10s(text string) []string {
	return match(text, ISBN10Regex)
}

// ParseVISACreditCards finds all VISA credit card numbers
func ParseVISACreditCards(text string) []string {
	return match(text, VISACreditCardRegex)
}

// ParseMCCreditCards finds all MasterCard credit card numbers
func ParseMCCreditCards(text string) []string {
	return match(text, MCCreditCardRegex)
}

// ParseMACAddresses finds all MAC addresses
func ParseMACAddresses(text string) []string {
	return match(text, MACAddressRegex)
}

// ParseIBANs finds all IBAN strings
func ParseIBANs(text string) []string {
	return match(text, IBANRegex)
}

// ParseGitRepos finds all git repository addresses which have protocol prefix
func ParseGitRepos(text string) []string {
	return match(text, GitRepoRegex)
}
