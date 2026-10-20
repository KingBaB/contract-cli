#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
VERSION_BEFORE="$(node -p "require('$ROOT_DIR/package.json').version")"
OUTPUT="$("$ROOT_DIR/scripts/release.sh" --version 0.1.3 --dry-run --skip-tests)"
VERSION_AFTER="$(node -p "require('$ROOT_DIR/package.json').version")"

if [ "$VERSION_BEFORE" != "$VERSION_AFTER" ]; then
  echo "dry-run changed package version: $VERSION_BEFORE -> $VERSION_AFTER" >&2
  exit 1
fi

printf '%s' "$OUTPUT" | grep -q "contract-cli-0.1.3-darwin-amd64.tar.gz"
printf '%s' "$OUTPUT" | grep -q "contract-cli-0.1.3-windows-arm64.zip"
printf '%s' "$OUTPUT" | grep -q "gh release create v0.1.3"
printf '%s' "$OUTPUT" | grep -q -- "--latest"
printf '%s' "$OUTPUT" | grep -q "npm publish --tag latest"

if printf '%s' "$OUTPUT" | grep -q -- "--prerelease"; then
  echo "stable release script should not mark release as prerelease" >&2
  exit 1
fi

"$ROOT_DIR/scripts/release.sh" --help | grep -q "Usage:"

echo "release script dry-run ok"
