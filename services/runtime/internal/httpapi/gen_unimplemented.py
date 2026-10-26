"""Sinh unimplemented.gen.go từ interface StrictServerInterface trong internal/api/api.gen.go."""
import re
import pathlib

here = pathlib.Path(__file__).parent
src = (here / "../api/api.gen.go").read_text()
block = src[src.index("type StrictServerInterface interface {"):]
block = block[:block.index("\n}\n")]
methods = re.findall(r"^\t(\w+)\(ctx context.Context, request (\w+)\) \((\w+), error\)$", block, re.M)
out = ["// Code generated from api.StrictServerInterface; DO NOT EDIT. Tạo lại bằng go generate ./internal/httpapi", "",
       "package httpapi", "", "import (", '\t"context"', "",
       '\t"github.com/NhatMinh0311/crucible/services/runtime/internal/api"',
       '\t"github.com/NhatMinh0311/crucible/services/runtime/internal/problem"', ")", "",
       "// Unimplemented trả 501 NOT_IMPLEMENTED cho mọi operation. Server nhúng kiểu này và chỉ ghi đè operation đã hiện thực.",
       "type Unimplemented struct{}", ""]
for name, req, resp in methods:
    out.append(f"func (Unimplemented) {name}(context.Context, api.{req}) (api.{resp}, error) {{\n\treturn nil, problem.ErrNotImplemented\n}}\n")
out.append("var _ api.StrictServerInterface = Unimplemented{}")
(here / "unimplemented.gen.go").write_text("\n".join(out) + "\n")
