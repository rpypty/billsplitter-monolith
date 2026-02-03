#!/usr/bin/env bash
set -euo pipefail

APP="billsplitter"
REPO_DIR="/root/Repos/billsplitter-monolith"
BIN_NAME="billsplitter"
BUILD_OUT="/tmp/${BIN_NAME}.new"
LIVE_BIN="/opt/billsplitter/${BIN_NAME}"
BACKUP_BIN="/opt/billsplitter/${BIN_NAME}.old"

echo "==> Pull latest code"
cd "$REPO_DIR"
git pull

echo "==> Build"
go mod download
go build -o "$BUILD_OUT" ./cmd/main

echo "==> Stop service"
systemctl stop "$APP"

echo "==> Swap binaries (atomic-ish)"
if [ -f "$LIVE_BIN" ]; then
  mv "$LIVE_BIN" "$BACKUP_BIN"
fi
mv "$BUILD_OUT" "$LIVE_BIN"
chown "$APP:$APP" "$LIVE_BIN"
chmod 755 "$LIVE_BIN"

echo "==> Start service"
systemctl start "$APP"

echo "==> Status"
systemctl --no-pager --full status "$APP" || true

echo "==> Last logs"
journalctl -u "$APP" -n 80 --no-pager
