#!/bin/sh
# Tạo CSDL và tài khoản riêng cho từng module (ADR-0001) và cho Keycloak.
# Chạy một lần khi khởi tạo volume PostgreSQL.
set -eu

create() {
  name="$1"; password="$2"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres <<SQL
CREATE ROLE "$name" LOGIN PASSWORD '$password';
CREATE DATABASE "$name" OWNER "$name";
REVOKE ALL ON DATABASE "$name" FROM PUBLIC;
SQL
}

create runtime  "$RUNTIME_DB_PASSWORD"
create content  "$CONTENT_DB_PASSWORD"
create learning "$LEARNING_DB_PASSWORD"
create keycloak "$KEYCLOAK_DB_PASSWORD"
