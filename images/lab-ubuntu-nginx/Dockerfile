# syntax=docker/dockerfile:1.10
# Image nền của bài lab: Ubuntu 24.04 có nginx và công cụ quản trị thường dùng.
# PID 1 là tini chạy sleep infinity; terminal, setup và lệnh kiểm tra đều đi qua exec (CR-MOD-M1 mục 5.3).

FROM ubuntu:24.04@sha256:68ee22004a9109b44973313b02ff7b1de250f1de0f5ecb628d73f579aaff33ac

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      nginx tini sudo cron \
      procps psmisc iproute2 net-tools iputils-ping iptables nftables dnsutils \
      curl ca-certificates less nano vim-tiny file lsof \
 && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8 TERM=xterm-256color
WORKDIR /root
LABEL org.opencontainers.image.title="crucible-labs/ubuntu-nginx" \
      org.opencontainers.image.description="Image nền bài lab Linux có nginx" \
      org.opencontainers.image.licenses="AGPL-3.0-only"
ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["sleep", "infinity"]
