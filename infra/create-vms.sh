#!/usr/bin/env bash
# Dựng mạng crucible-net, pool lưu trữ và hai máy ảo cr-server, cr-lab (CR-DOC-04 mục 3; ADR-0009).
# Cấu hình bên trong máy ảo do cloud-init đảm nhận (infra/cloud-init/). Chạy lại an toàn: bỏ qua phần đã có.
#
# Biến môi trường (giá trị mặc định theo CR-DOC-04 mục 3.2):
#   CRUCIBLE_POOL_DIR      thư mục pool lưu trữ (mặc định $HOME/crucible-vms; qemu cần quyền đi qua $HOME)
#   CRUCIBLE_SECRETS_DIR   nơi giữ khóa SSH của runtime-api (mặc định $HOME/.config/crucible)
#   ADMIN_SSH_PUBKEY       khóa công khai của quản trị viên (mặc định ~/.ssh/id_ed25519.pub hoặc id_rsa.pub)
#   SERVER_VCPUS, SERVER_MEMORY_MIB, LAB_VCPUS, LAB_MEMORY_MIB
set -euo pipefail

export LIBVIRT_DEFAULT_URI=qemu:///system
HERE="$(cd "$(dirname "$0")" && pwd)"

POOL=crucible
POOL_DIR="${CRUCIBLE_POOL_DIR:-$HOME/crucible-vms}"
SECRETS_DIR="${CRUCIBLE_SECRETS_DIR:-$HOME/.config/crucible}"
CACHE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/crucible"
IMAGE_URL=https://cloud-images.ubuntu.com/noble/current
IMAGE_NAME=noble-server-cloudimg-amd64.img
BASE_VOL=noble-base.img

SERVER_VCPUS="${SERVER_VCPUS:-4}"
SERVER_MEMORY_MIB="${SERVER_MEMORY_MIB:-10240}"
LAB_VCPUS="${LAB_VCPUS:-8}"
LAB_MEMORY_MIB="${LAB_MEMORY_MIB:-16384}"

log() { printf '==> %s\n' "$*"; }

admin_pubkey() {
  local f
  for f in "${ADMIN_SSH_PUBKEY:-}" ~/.ssh/id_ed25519.pub ~/.ssh/id_rsa.pub; do
    [ -n "$f" ] && [ -f "$f" ] && { cat "$f"; return; }
  done
  echo "Không tìm thấy khóa SSH công khai; đặt ADMIN_SSH_PUBKEY" >&2
  exit 1
}

# 1. Mạng NAT 10.77.0.0/24, địa chỉ tĩnh theo MAC
if ! virsh net-info crucible-net >/dev/null 2>&1; then
  log "Tạo mạng crucible-net"
  virsh net-define "$HERE/libvirt/crucible-net.xml"
  virsh net-autostart crucible-net
fi
virsh net-info crucible-net | grep -q 'Active: *yes' || virsh net-start crucible-net

# 2. Pool lưu trữ dạng thư mục; libvirtd (root) tạo thư mục
if ! virsh pool-info "$POOL" >/dev/null 2>&1; then
  log "Tạo pool $POOL tại $POOL_DIR"
  virsh pool-define-as "$POOL" dir --target "$POOL_DIR"
  virsh pool-build "$POOL"
  virsh pool-autostart "$POOL"
fi
virsh pool-info "$POOL" | grep -q 'State: *running' || virsh pool-start "$POOL"

# 3. Image gốc Ubuntu Server 24.04 cloud, kiểm tra SHA-256 rồi đưa vào pool
if ! virsh vol-info --pool "$POOL" "$BASE_VOL" >/dev/null 2>&1; then
  mkdir -p "$CACHE_DIR"
  curl -fsSL "$IMAGE_URL/SHA256SUMS" -o "$CACHE_DIR/SHA256SUMS"
  if ! (cd "$CACHE_DIR" && grep " \*$IMAGE_NAME\$" SHA256SUMS | sha256sum -c --status); then
    log "Tải $IMAGE_NAME"
    curl -fL "$IMAGE_URL/$IMAGE_NAME" -o "$CACHE_DIR/$IMAGE_NAME"
    (cd "$CACHE_DIR" && grep " \*$IMAGE_NAME\$" SHA256SUMS | sha256sum -c)
  fi
  log "Đưa image gốc vào pool"
  virsh vol-create-as "$POOL" "$BASE_VOL" "$(stat -c %s "$CACHE_DIR/$IMAGE_NAME")" --format raw
  virsh vol-upload --pool "$POOL" "$BASE_VOL" "$CACHE_DIR/$IMAGE_NAME"
fi
BASE_PATH="$(virsh vol-path --pool "$POOL" "$BASE_VOL")"

# 4. Khóa SSH của runtime-api: chỉ nằm trên máy quản trị và cr-server, không vào kho mã
mkdir -p "$SECRETS_DIR" && chmod 700 "$SECRETS_DIR"
[ -f "$SECRETS_DIR/runtime-api" ] || ssh-keygen -q -t ed25519 -N '' -C runtime-api -f "$SECRETS_DIR/runtime-api"

# 5. Sinh user-data từ mẫu vào thư mục tạm (chứa khóa riêng nên xóa sau khi dùng)
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
ADMIN_KEY="$(admin_pubkey)"
RUNTIME_PUB="$(cat "$SECRETS_DIR/runtime-api.pub")"
RUNTIME_KEY_B64="$(base64 -w0 "$SECRETS_DIR/runtime-api")"
for name in cr-server cr-lab; do
  sed -e "s|@ADMIN_SSH_PUBKEY@|$ADMIN_KEY|" \
      -e "s|@RUNTIME_API_PUBKEY@|$RUNTIME_PUB|" \
      -e "s|@RUNTIME_API_KEY_B64@|$RUNTIME_KEY_B64|" \
      "$HERE/cloud-init/$name.yaml" > "$WORK/$name.yaml"
done

# 6. Máy ảo: đĩa hệ điều hành là lớp qcow2 trên image gốc; cr-lab thêm đĩa riêng cho /var/lib/docker
create_vm() {
  local name=$1 vcpus=$2 mem=$3 mac=$4; shift 4
  if virsh dominfo "$name" >/dev/null 2>&1; then
    log "$name đã tồn tại, bỏ qua"
    return
  fi
  log "Tạo máy ảo $name ($vcpus vCPU, $mem MiB)"
  virt-install --name "$name" --vcpus "$vcpus" --memory "$mem" \
    --osinfo linux2022 --import \
    --disk "pool=$POOL,size=$1,backing_store=$BASE_PATH,backing_format=qcow2,bus=virtio" \
    "${@:2}" \
    --network "network=crucible-net,mac=$mac,model=virtio" \
    --cloud-init "user-data=$WORK/$name.yaml" \
    --graphics none --noautoconsole
}
create_vm cr-server "$SERVER_VCPUS" "$SERVER_MEMORY_MIB" 52:54:00:77:00:10 60
create_vm cr-lab "$LAB_VCPUS" "$LAB_MEMORY_MIB" 52:54:00:77:00:11 30 \
  --disk "pool=$POOL,size=40,format=qcow2,bus=virtio"

# 7. Chờ cloud-init chạy xong trên cả hai máy
for ip in 10.77.0.10 10.77.0.11; do
  ssh-keygen -R "$ip" >/dev/null 2>&1 || true
  log "Chờ cloud-init tại $ip"
  until ssh -o StrictHostKeyChecking=accept-new -o ConnectTimeout=3 -o BatchMode=yes "ops@$ip" true 2>/dev/null; do
    sleep 5
  done
  ssh "ops@$ip" cloud-init status --wait --long | tail -n 3
done
