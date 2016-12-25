package util

import (
	"bytes"
	"html/template"
	"strings"

	"golang.org/x/net/html"
)

// StripTags 过滤HTML标签，保留指定的tag
func StripTags(s []byte, tags ...string) []byte {
	r := bytes.NewBuffer(nil)
	_, err := r.Write(s)
	if err != nil {
		return s
	}

	re := make([]byte, 0, len(s))
	z := html.NewTokenizer(r)
	var t html.TokenType
	for {
		t = z.Next()
		if t == html.ErrorToken {
			break
		}
		if t == html.TextToken {
			for _, v := range z.Raw() {
				if v == ' ' || v == '\t' {
					continue
				}
				re = append(re, v)
			}
			continue
		}
		tagName, _ := z.TagName()
		if tags != nil && InSlice(tags, string(tagName), "string") {
			re = append(re, z.Raw()...)
		}
	}

	return re
}

// Nl2br 将换行转换成HTML中的BR
func Nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
}
