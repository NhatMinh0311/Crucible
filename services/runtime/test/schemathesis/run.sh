#!/usr/bin/env bash
# Chạy Schemathesis với các operation của chế độ docker. Cần runtime-api và Keycloak của dev/compose.yaml đang chạy.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../../.." && pwd)"
KC=http://localhost:8180/realms/crucible/protocol/openid-connect/token
TOKEN=$(curl -sf "$KC" -d grant_type=client_credentials -d client_id=content-service \
  -d client_secret=content-service-dev-secret | sed -E 's/.*"access_token":"([^"]+)".*/\1/')
docker run --rm --network host \
  -v "$ROOT/contracts:/contracts:ro" -v "$ROOT/services/runtime/test/schemathesis:/cfg:ro" \
  schemathesis/schemathesis:4.31.1 --config-file /cfg/schemathesis.toml run /contracts/runtime-api.yaml \
  --url http://localhost:8081 -H "Authorization: Bearer $TOKEN" \
  --include-operation-id createSession --include-operation-id getSession --include-operation-id extendSession \
  --include-operation-id terminateSession --include-operation-id execInSession \
  --max-examples "${MAX_EXAMPLES:-25}" "$@"
