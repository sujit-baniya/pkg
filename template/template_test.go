package template

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func BenchmarkLoopReplace(b *testing.B) {
	template := "https://{{host}}/?q={{query}}&foo={{bar}}{{bar}}"
	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{
			"host":  "google.com",
			"query": url.QueryEscape("hello=world"),
			"bar":   "foobar",
		}
		for k, v := range data {
			template = strings.ReplaceAll(template, "{{"+k+"}}", fmt.Sprintf("%v", v))
		}
	}
}

func BenchmarkTemplate_Parse(b *testing.B) {
	tmpl := New("https://{{host}}/?q={{query}}&foo={{bar}}{{bar}}", "", "")
	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{
			"host":  "google.com",
			"query": url.QueryEscape("hello=world"),
			"bar":   "foobar",
		}
		tmpl.Parse(data)
	}
}
