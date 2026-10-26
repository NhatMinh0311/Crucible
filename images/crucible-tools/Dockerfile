# syntax=docker/dockerfile:1.10
# Bộ công cụ chỉ đọc cho lệnh kiểm tra, gắn vào môi trường tại /opt/crucible/bin (ADR-0011, CR-DOC-03 mục 7).
# Mọi chương trình liên kết tĩnh để chạy được trên image bài lab bất kỳ. Phiên bản và mã băm được ghim.

ARG BUSYBOX_IMAGE=busybox:1.37.0-musl@sha256:5cec3fc171c87218698e85a52af7087de727372aae264a787b8112901a5b0092
ARG JQ_VERSION=1.8.2
ARG CURL_VERSION=8.23.0

FROM ${BUSYBOX_IMAGE} AS busybox

FROM busybox AS fetch-amd64
ARG JQ_VERSION CURL_VERSION
ADD --checksum=sha256:b1c22172dd303f3be49e935aa56aa48a8b7a46e0bc838b4997d3bb451495870f \
    https://github.com/jqlang/jq/releases/download/jq-${JQ_VERSION}/jq-linux-amd64 /fetch/jq
ADD --checksum=sha256:049f982cb7db34e745cc0be0bdce2246b51fbf7ef17fc95ce37980f711ca91bb \
    https://github.com/stunnel/static-curl/releases/download/${CURL_VERSION}/curl-linux-x86_64-musl-${CURL_VERSION}.tar.xz /fetch/curl.tar.xz

FROM busybox AS fetch-arm64
ARG JQ_VERSION CURL_VERSION
ADD --checksum=sha256:8b85c817833814ddca00a144c33705546355afccf0cf39b188f3cdb48b852309 \
    https://github.com/jqlang/jq/releases/download/jq-${JQ_VERSION}/jq-linux-arm64 /fetch/jq
ADD --checksum=sha256:4db6c9b5e7f61f6a8cfa7f83e8c24e54259334bbb46743f6207ffb67bc0d1c3a \
    https://github.com/stunnel/static-curl/releases/download/${CURL_VERSION}/curl-linux-aarch64-musl-${CURL_VERSION}.tar.xz /fetch/curl.tar.xz

FROM fetch-${TARGETARCH} AS assemble
RUN mkdir -p /opt/crucible/bin \
 && cp /bin/busybox /opt/crucible/bin/busybox \
 && /opt/crucible/bin/busybox --install -s /opt/crucible/bin \
 && tar -xJf /fetch/curl.tar.xz -C /fetch \
 && install -m 0755 /fetch/curl /opt/crucible/bin/curl \
 && install -m 0755 /fetch/jq /opt/crucible/bin/jq

FROM scratch
COPY --from=assemble /opt/crucible/bin /opt/crucible/bin
ENV PATH=/opt/crucible/bin
LABEL org.opencontainers.image.title="crucible-tools" \
      org.opencontainers.image.description="Bộ công cụ chỉ đọc cho lệnh kiểm tra của Crucible" \
      org.opencontainers.image.licenses="AGPL-3.0-only"
CMD ["/opt/crucible/bin/sh"]
