#!/usr/bin/env bash
# Xóa hai máy ảo cr-server, cr-lab cùng đĩa của chúng; giữ mạng, pool và image gốc để dựng lại bằng create-vms.sh.
# Dữ liệu trên máy ảo mất hẳn: chỉ dùng trên môi trường thử hoặc sau khi đã sao lưu (CR-DOC-04 mục 8).
set -euo pipefail

export LIBVIRT_DEFAULT_URI=qemu:///system

names=("$@")
[ $# -gt 0 ] || names=(cr-server cr-lab)

for name in "${names[@]}"; do
  virsh dominfo "$name" >/dev/null 2>&1 || continue
  virsh destroy "$name" >/dev/null 2>&1 || true
  virsh undefine "$name" --remove-all-storage
done
