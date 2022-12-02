package template

import (
	"fmt"
	"strings"
)

type Template struct {
	Content     string `json:"content"`
	ByteContent []byte `json:"byte_content"`
	Prefix      string `json:"prefix"`
	Suffix      string `json:"suffix"`
}

func (t *Template) Parse(data map[string]interface{}) string {
	content := t.Content
	for k, v := range data {
		content = strings.ReplaceAll(content, t.Prefix+k+t.Suffix, fmt.Sprintf("%v", v))
	}
	return content
}

func New(content string, prefix, suffix string) *Template {
	if prefix == "" {
		prefix = "{{"
	}
	if suffix == "" && prefix == "{{" {
		suffix = "}}"
	}
	return &Template{Content: content, Prefix: prefix, Suffix: suffix}
}
