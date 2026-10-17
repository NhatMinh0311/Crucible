#!/usr/bin/env bash
# Kiểm tra dịch vụ phụ thuộc của môi trường phát triển và tạo bucket còn thiếu.
# Chạy sau: docker compose -f dev/compose.yaml up -d --wait
set -euo pipefail
cd "$(dirname "$0")"
set -a; source .env; set +a
# Tìm container theo tên dịch vụ, dùng được cả khi chạy trên máy và trong Dev Container.
ctr() { docker ps --filter "label=com.docker.compose.service=$1" --format '{{.Names}}' | head -1; }
ok()   { printf '  \033[32m✔\033[0m %s\n' "$1"; }
fail() { printf '  \033[31m✘\033[0m %s\n' "$1"; exit 1; }

echo "PostgreSQL"
for db in runtime content learning; do
  var="$(echo "$db" | tr a-z A-Z)_DB_PASSWORD"
  docker exec -e PGPASSWORD="${!var}" "$(ctr postgres)" psql -h 127.0.0.1 -U "$db" -d "$db" -tAc 'select 1' >/dev/null \
    && ok "CSDL $db, tài khoản $db" || fail "CSDL $db"
done

echo "RabbitMQ"
curl -sf -u "$RABBITMQ_USER:$RABBITMQ_PASSWORD" "http://localhost:${RABBITMQ_UI_PORT}/api/overview" >/dev/null \
  && ok "AMQP localhost:${AMQP_PORT}, giao diện http://localhost:${RABBITMQ_UI_PORT}" || fail "RabbitMQ"

echo "SeaweedFS (S3)"
for b in recordings lab-assets; do
  echo "s3.bucket.create -name $b" | docker exec -i "$(ctr seaweedfs)" weed shell >/dev/null 2>&1 || true
done
buckets=$(echo "s3.bucket.list" | docker exec -i "$(ctr seaweedfs)" weed shell 2>/dev/null)
for b in recordings lab-assets; do
  grep -q "$b" <<<"$buckets" && ok "bucket $b tại http://localhost:${S3_PORT}" || fail "bucket $b"
done

echo "Keycloak"
KC="http://localhost:${KEYCLOAK_PORT}/realms/crucible/protocol/openid-connect/token"
claims() { python3 -c "import sys,json,base64;t=json.load(sys.stdin)['access_token'].split('.')[1];t+='='*(-len(t)%4);print(base64.urlsafe_b64decode(t).decode())"; }
svc=$(curl -sf -d grant_type=client_credentials -d client_id=content-service -d client_secret=content-service-dev-secret "$KC" | claims)
python3 -c "import json,sys;c=json.loads(sys.argv[1]);assert 'runtime.exec' in c['scope'].split() and 'runtime-api' in c['aud']" "$svc" \
  && ok "content-service: client credentials có scope runtime.exec, audience runtime-api" || fail "token dịch vụ"
usr=$(curl -sf -d grant_type=password -d client_id=crucible-web -d username=sv1 -d password=crucible -d scope=openid "$KC" | claims)
python3 -c "import json,sys;c=json.loads(sys.argv[1]);assert 'STUDENT' in c['realm_access']['roles'] and {'runtime-api','content-service','learning-service'}<=set(c['aud'])" "$usr" \
  && ok "sv1: vai trò STUDENT, audience ba dịch vụ" || fail "token người dùng"

echo "Mailpit"
curl -sf "http://localhost:${MAILPIT_UI_PORT}/api/v1/info" >/dev/null \
  && ok "SMTP localhost:${SMTP_PORT}, giao diện http://localhost:${MAILPIT_UI_PORT}" || fail "Mailpit"

echo "Môi trường phát triển sẵn sàng."
