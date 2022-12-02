package permission

import "strings"

type validationRule int

const (
	matchAll validationRule = iota
	atLeastOne
)

// MatchAll is an option that defines all permissions
// or roles should match the user.
var MatchAll = func(o *Options) {
	o.ValidationRule = matchAll
}

// AtLeastOne is an option that defines at least on of
// permissions or roles should match to pass.
var AtLeastOne = func(o *Options) {
	o.ValidationRule = atLeastOne
}

// ParserFunc is used for parsing the permission
// to extract object and action usually
type ParserFunc func(str string) []string

func permissionParserWithSeparator(sep string) ParserFunc {
	return func(str string) []string {
		return strings.Split(str, sep)
	}
}

// ParserWithSeparator is an option that parses permission
// with separators
func ParserWithSeparator(sep string) func(o *Options) {
	return func(o *Options) {
		o.PermissionParser = permissionParserWithSeparator(sep)
	}
}

// Options holds options of middleware
type Options struct {
	PermissionParser ParserFunc
	ValidationRule   validationRule
}
