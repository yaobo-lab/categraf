#!/usr/bin/env sh
set -eu

URL="${1:-${HEARTBEAT_URL:-http://127.0.0.1:17000/v1/n9e/heartbeat}}"
HOST_IP="${HOST_IP:-127.0.0.1}"

curl -i -sS -X POST "$URL" \
	-H "Content-Type: application/json" \
	-H "User-Agent: categraf/$HOST_IP"
