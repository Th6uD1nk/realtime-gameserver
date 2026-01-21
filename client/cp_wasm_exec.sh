#!/bin/bash

echo "Getting Go version..."
GO_FULL_VER=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR_MINOR=$(echo $GO_FULL_VER | cut -d. -f1-2)
echo "Detected Go version: $GO_FULL_VER (major.minor: $GO_MAJOR_MINOR)"

URLS=(
  "https://raw.githubusercontent.com/golang/go/go$GO_FULL_VER/lib/wasm/wasm_exec.js"
  "https://raw.githubusercontent.com/golang/go/go$GO_MAJOR_MINOR/lib/wasm/wasm_exec.js"
  "https://raw.githubusercontent.com/golang/go/master/lib/wasm/wasm_exec.js"
)

SUCCESS=0
for url in "${URLS[@]}"; do
  echo "Trying to download: $url"
  curl -f -# -o ./wasm_exec.js "$url"
  if [ $? -eq 0 ]; then
    echo "Downloaded wasm_exec.js from $url"
    SUCCESS=1
    break
  fi
done

if [ $SUCCESS -eq 0 ]; then
  echo "ERROR: Could not download wasm_exec.js from any known location"
  exit 1
fi

