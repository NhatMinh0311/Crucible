#!/usr/bin/env bash
# Chạy một lần sau khi Dev Container được tạo.
set -euo pipefail
pip install --user --quiet jsonschema==4.23.0 pyyaml==6.0.2
echo "Phiên bản công cụ:"
java -version 2>&1 | head -1
node --version
go version
python3 --version
docker compose version
