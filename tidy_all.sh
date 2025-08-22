#!/usr/bin/env bash
set -euo pipefail

find . -name "go.mod" -execdir bash -c '
    echo "Tidying module in $(pwd)"
    go mod tidy
' \;

echo "✅ All modules tidied!"