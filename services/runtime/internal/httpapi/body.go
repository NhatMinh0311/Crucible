package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"
)

// maxBody giới hạn kích thước yêu cầu; đặc tả môi trường tối đa 256 KiB (RenderedEnvironment).
const maxBody = 1 << 20

// normalizeJSON giới hạn kích thước thân yêu cầu và viết lại số nguyên ở dạng thập phân (10.0, 1e1) thành 10. JSON
// Schema coi các giá trị này là số nguyên nhưng encoding/json của Go không đọc được vào kiểu int.
func normalizeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			next.ServeHTTP(w, r)
			return
		}
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
		if err != nil {
			problem.Write(w, problem.Validation("thân yêu cầu vượt 1 MiB hoặc không đọc được"))
			return
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber()
		var v any
		if dec.Decode(&v) == nil {
			if changed := integralNumbers(&v); changed {
				if b, err := json.Marshal(v); err == nil {
					raw = b
				}
			}
		}
		r.Body = io.NopCloser(bytes.NewReader(raw))
		r.ContentLength = int64(len(raw))
		next.ServeHTTP(w, r)
	})
}

func integralNumbers(v *any) bool {
	changed := false
	switch t := (*v).(type) {
	case map[string]any:
		for k, e := range t {
			if integralNumbers(&e) {
				t[k], changed = e, true
			}
		}
	case []any:
		for i := range t {
			if integralNumbers(&t[i]) {
				changed = true
			}
		}
	case json.Number:
		s := string(t)
		if !strings.ContainsAny(s, ".eE") {
			return false
		}
		f, ok := new(big.Float).SetString(s)
		if !ok || !f.IsInt() {
			return false
		}
		i, _ := f.Int(nil)
		*v = json.Number(i.String())
		return true
	}
	return changed
}
