package httpapi

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIntegralNumbers(t *testing.T) {
	for in, want := range map[string]string{
		`{"minutes":10.0}`:                      `{"minutes":10}`,
		`{"a":[1e1,2.5,{"b":3.00}],"c":"10.0"}`: `{"a":[10,2.5,{"b":3}],"c":"10.0"}`,
		`{"minutes":15}`:                        `{"minutes":15}`,
	} {
		var v any
		dec := json.NewDecoder(stringsReader(in))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			t.Fatal(err)
		}
		integralNumbers(&v)
		got, _ := json.Marshal(v)
		if string(got) != want {
			t.Errorf("%s: mong đợi %s, nhận %s", in, want, got)
		}
	}
}

func stringsReader(s string) *strings.Reader { return strings.NewReader(s) }
