# syntax=docker/dockerfile:1.10
# Image runtime-api cho môi trường phát triển (dev/compose.yaml, profile runtime). Chế độ docker cần socket của Docker
# Engine trên máy nên image chạy bằng root; không dùng image này trên cụm.
FROM golang:1.26.5-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/runtime-api ./cmd/runtime-api

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/runtime-api /runtime-api
EXPOSE 8081
ENTRYPOINT ["/runtime-api"]
