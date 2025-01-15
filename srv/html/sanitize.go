package htmlsrv

import (
	"github.com/microcosm-cc/bluemonday"
)

var policy *bluemonday.Policy = bluemonday.UGCPolicy().AllowElements(
	"b",
	"blockquote",
	"code",
	"dd",
	"del",
	"dl",
	"dt",
	"em",
	"i",
	"ins",
	"kbd",
	"li",
	"mark",
	"ol",
	"p",
	"pre",
	"q",
	"samp",
	"small",
	"strong",
	"sub",
	"sup",
	"ul",
)

func init() {
	if nil == policy {
		panic("nil bluemonday HTML sanitization policy")
	}
}

func SanitizeBytes(value []byte) []byte {
	return policy.SanitizeBytes(value)
}

func SanitizeString(value string) string {
	return policy.Sanitize(value)
}
