package email

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"time"
)

type Email struct {
	Email      string `json:"email"`
	Domain     string
	Disposable bool   `json:"is_disposable"`
	ValidateMX bool   `json:"validate_mx"`
	Free       bool   `json:"is_free"`
	MxError    string `json:"mx_error,omitempty"`
	NsError    string `json:"ns_error,omitempty"`
	HostError  string `json:"host_error,omitempty"`
	IpError    string `json:"ip_error,omitempty"`
	Error      string `json:"error,omitempty"`
	Valid      bool   `json:"is_valid"`
	mx         []*net.MX
}

type EmailList struct {
	Emails           []string `json:"emails"`
	RemoveDisposable bool     `json:"remove_disposable"`
}

type Emails struct {
	Emails []Email `json:"emails"`
}

const forceDisconnectAfter = time.Second * 50

const (
	emptyString string = ""
)

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// IsValid Validate - validates an email address via all options
func (e *Email) IsValid() {
	e.Valid = true
	e.ValidateFormat()
	e.IsDisposable()
	if e.ValidateMX {
		e.ValidateDomainRecords()
	}
	// e.ValidateHostAndUser("smtp-relay.sendinblue.com", "info@verishore.com", e.mx)
}

// Check Validate - validates an email address via all options
func (e *Email) Check() {
	e.Valid = true
	e.ValidateFormat()
	e.IsDisposable()
	e.IsFree()
	// e.ValidateDomainRecords()
	// e.ValidateHostAndUser("smtp-relay.sendinblue.com", "info@verishore.com", e.mx)
}

// IsValid Validate - validates an email address via all options
func (e *EmailList) IsValid() Emails {
	emails := Emails{}
	email := Email{}
	for _, val := range e.Emails {
		email.Email = val
		email.IsValid()
		emails.Emails = append(emails.Emails, email)
	}
	return emails
}

// Clean - validates an email address via all options
func (e *EmailList) Clean() Emails {
	emails := Emails{}
	email := Email{}
	for _, val := range e.Emails {
		email.Email = val
		email.IsValid()
		if email.Valid {
			if e.RemoveDisposable && !email.Disposable {
				emails.Emails = append(emails.Emails, email)
			} else if !e.RemoveDisposable {
				emails.Emails = append(emails.Emails, email)
			}
		}
	}
	return emails
}

// Stats - validates an email address via all options
func (e *EmailList) Stats() map[string]int {
	email := Email{}
	totalCount := len(e.Emails)
	invalidCount := 0
	disposableCount := 0
	freeCount := 0
	for _, val := range e.Emails {
		email.Email = val
		email.IsValid()
		if !email.Valid {
			invalidCount++
		}
		if email.Disposable {
			disposableCount++
		}
		if email.Free {
			freeCount++
		}
	}
	mp := map[string]int{
		"total_count":      totalCount,
		"invalid_count":    invalidCount,
		"disposable_count": disposableCount,
		"free_count":       freeCount,
	}
	return mp
}

// ValidateFormat - validates an email address meets rfc 822 format via a regex
func (e *Email) ValidateFormat() {
	_, domain, err := validateFormatAndSplit(e.Email)
	if err != nil {
		e.Valid = false
		e.Error = err.Error()
	}
	e.Domain = domain
}

// ValidateDomainRecords - validates an email address domain's NS and MX records via a DNS lookup
func (e *Email) ValidateDomainRecords() {
	// Added NS check as some ISPs hijack the MX record lookup :(
	nsRecords, err := net.LookupNS(e.Domain)
	if err != nil || len(nsRecords) == 0 {
		e.Valid = false
		e.NsError = "Invalid email domain, unable to find Name Servers records"
		return
	}
	mx, err := net.LookupMX(e.Domain)
	if err != nil {
		e.Valid = false
		e.MxError = "Invalid email domain no MX records found"
		return
	}
	e.mx = mx
	// e.ValidateHostAndUser("smtp.google.com", "s.baniya.np@gmail.com", mx)
	if _, err := net.LookupIP(e.Domain); err != nil {
		e.Valid = false
		e.IpError = "Invalid email domain no IP records found"
		return
	}
}

// Normalize - Trim whitespaces, extra dots in the hostname and converts to Lowercase.
func Normalize(email string) string {

	email = strings.TrimSpace(email)
	email = strings.TrimRight(email, ".")
	email = strings.ToLower(email)

	return email
}

func validateFormatAndSplit(email string) (username string, domain string, err error) {
	if len(email) < 6 || len(email) > 254 {
		return emptyString, emptyString, errors.New("Invalid Email Format")
	}

	// Regex matches as per rfc 822 https://tools.ietf.org/html/rfc822
	if !emailRegexp.MatchString(email) {
		return emptyString, emptyString, errors.New("Invalid Email Format")
	}

	i := strings.LastIndexByte(email, '@')
	username = email[:i]
	domain = email[i+1:]

	if len(username) > 64 {
		return emptyString, emptyString, errors.New("Invalid Email Format")
	}

	return username, domain, nil
}

func GetDomainOfEmail(email string) string {
	i := strings.LastIndexByte(email, '@')
	return email[i+1:]
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (e *Email) IsDisposable() {
	e.Disposable = isDisposable(e.Domain)
}

func (e *Email) IsFree() {
	e.Free = IsFreeDomain(e.Domain)
}

func ValidateEmail(email string, validateMX ...bool) Email {
	e := Email{Email: email}
	if len(validateMX) > 0 {
		e.ValidateMX = validateMX[0]
	}
	e.IsValid()
	e.IsFree()
	return e
}

func ValidateEmailList(emails []string, removeDisposable bool) EmailList {
	e := EmailList{Emails: emails, RemoveDisposable: removeDisposable}
	e.IsValid()
	return e
}

func EmailStats(emails []string) map[string]int {
	e := EmailList{Emails: emails}
	return e.Stats()
}

func CleanEmailList(emails []string) Emails {
	e := EmailList{Emails: emails}
	return e.Clean()
}
