package httpapi

import (
	"regexp"
	"slices"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// allowIndex tìm các phương thức của một đường dẫn để trả header Allow cùng 405 (RFC 9110 mục 15.5.6).
type allowIndex []struct {
	re      *regexp.Regexp
	methods string
}

var pathParam = regexp.MustCompile(`\\\{[^}]+\\\}`)

func newAllowIndex(spec *openapi3.T) allowIndex {
	var idx allowIndex
	for tmpl, item := range spec.Paths.Map() {
		re := regexp.MustCompile("^" + pathParam.ReplaceAllString(regexp.QuoteMeta(tmpl), "[^/]+") + "$")
		var ms []string
		for m := range item.Operations() {
			ms = append(ms, strings.ToUpper(m))
		}
		slices.Sort(ms)
		idx = append(idx, struct {
			re      *regexp.Regexp
			methods string
		}{re, strings.Join(ms, ", ")})
	}
	return idx
}

func (a allowIndex) methods(path string) string {
	for _, e := range a {
		if e.re.MatchString(path) {
			return e.methods
		}
	}
	return ""
}
