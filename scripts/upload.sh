#!/usr/bin/env bash
set -eu

API_BASE=${1:-http://localhost:8000}
FILE=${2:-}

API=${API_BASE%/}

if [ -z "$FILE" ]; then
  echo "Usage: $0 [API_BASE_URL] /path/to/image.jpg"
  echo "Example: $0 http://localhost:8000 ./sample.jpg"
  exit 2
fi

resp=$(curl -s -w "\n%{http_code}" -X POST "$API/api/v1/uploads" -F "file=@${FILE}")
http=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')

if [ "$http" != "201" ]; then
  echo "Upload failed (HTTP $http)"
  echo "$body"
  exit 3
fi

# pretty print URL or raw JSON
url=$(echo "$body" | sed -n 's/.*"url"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
if [ -n "$url" ]; then
  echo "Uploaded: $url"
else
  echo "$body"
fi
