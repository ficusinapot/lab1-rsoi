#!/bin/sh
set -eu

CONFIG_PATH="${RSOI_CONFIG_PATH:-/app/configs/config.yaml}"

if [ -n "${DATABASE_URL:-}" ]; then
    atlas migrate apply \
        --dir "file:///app/migrations" \
        --url "$DATABASE_URL"
fi

exec /app/service -c "$CONFIG_PATH"
