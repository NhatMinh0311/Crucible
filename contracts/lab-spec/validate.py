#!/usr/bin/env python3
"""Kiểm tra đặc tả bài lab: lược đồ JSON Schema và các quy tắc ngữ nghĩa tại CR-DOC-03 mục 4.

Cách dùng: python3 validate.py <tệp.yaml> [<tệp.yaml> ...]
Mã thoát 0 khi mọi tệp hợp lệ.
"""
import json
import re
import sys
from pathlib import Path

import yaml
from jsonschema import Draft202012Validator

SCHEMA_PATH = Path(__file__).with_name("lab-spec.schema.json")
PLACEHOLDER = re.compile(r"\{\{\s*([^}\s]*)\s*\}\}")


def templated_fields(spec):
    """Liệt kê (vị trí, chuỗi) của các trường được kết xuất tham số."""
    env = spec.get("environment", {})
    for i, f in enumerate(env.get("files", [])):
        yield f"environment.files[{i}].content", f["content"]
    if "setup" in env:
        yield "environment.setup", env["setup"]
    for i, fault in enumerate(spec.get("faults", {}).get("from", [])):
        yield f"faults.from[{i}].inject", fault["inject"]
    for i, step in enumerate(spec["steps"]):
        yield f"steps[{i}].description", step["description"]
        yield f"steps[{i}].check.run", step["check"]["run"]
        if "failureMessage" in step["check"]:
            yield f"steps[{i}].check.failureMessage", step["check"]["failureMessage"]
        for j, hint in enumerate(step.get("hints", [])):
            yield f"steps[{i}].hints[{j}].text", hint["text"]
    for i, sol in enumerate(spec["solution"]):
        yield f"solution[{i}].run", sol["run"]


def semantic_errors(spec):
    errors = []
    params = spec.get("params", {})
    used = set()

    for where, text in templated_fields(spec):
        for name in PLACEHOLDER.findall(text):
            if name not in params:
                errors.append(f"{where}: tham số '{name}' chưa được khai báo")
            used.add(name)
    for name in sorted(set(params) - used):
        errors.append(f"params.{name}: tham số không được sử dụng")

    for name, gen in params.items():
        if "range" in gen and gen["range"]["min"] >= gen["range"]["max"]:
            errors.append(f"params.{name}.range: min phải nhỏ hơn max")

    for i, group in enumerate(spec.get("constraints", {}).get("distinct", [])):
        for name in group:
            if name not in params:
                errors.append(f"constraints.distinct[{i}]: tham số '{name}' chưa được khai báo")

    step_ids = [s["id"] for s in spec["steps"]]
    if len(step_ids) != len(set(step_ids)):
        errors.append("steps: mã bước bị trùng")
    solved = [s["step"] for s in spec["solution"]]
    if sorted(solved) != sorted(set(step_ids)):
        errors.append("solution: mỗi bước phải có đúng một phần lời giải và chỉ tham chiếu bước tồn tại")

    if sum(s["points"] for s in spec["steps"]) == 0:
        errors.append("steps: tổng điểm phải lớn hơn 0")
    for i, step in enumerate(spec["steps"]):
        if sum(h["penalty"] for h in step.get("hints", [])) > 1:
            errors.append(f"steps[{i}].hints: tổng tỉ lệ trừ điểm vượt quá 1")
    if all(s.get("expectInitially", "fail") == "pass" for s in spec["steps"]):
        errors.append("steps: ít nhất một bước phải không đạt trên môi trường mới dựng")

    faults = spec.get("faults")
    if faults:
        ids = [f["id"] for f in faults["from"]]
        if len(ids) != len(set(ids)):
            errors.append("faults.from: mã lỗi cài sẵn bị trùng")
        if faults["pick"] >= len(ids):
            errors.append("faults.pick: phải nhỏ hơn số lỗi trong danh sách")

    return errors


def main(paths):
    validator = Draft202012Validator(json.loads(SCHEMA_PATH.read_text(encoding="utf-8")))
    failed = False
    for path in paths:
        spec = yaml.safe_load(Path(path).read_text(encoding="utf-8"))
        errors = [
            f"{'/'.join(map(str, e.absolute_path)) or '(gốc)'}: {e.message}"
            for e in sorted(validator.iter_errors(spec), key=lambda e: list(e.absolute_path))
        ]
        if not errors:
            errors = semantic_errors(spec)
        if errors:
            failed = True
            print(f"KHÔNG HỢP LỆ  {path}")
            for err in errors:
                print(f"  - {err}")
        else:
            print(f"Hợp lệ       {path}")
    return 1 if failed else 0


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(2)
    sys.exit(main(sys.argv[1:]))
