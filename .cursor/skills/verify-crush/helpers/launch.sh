#!/usr/bin/env bash
set -euo pipefail

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
TEST_DIR="/tmp/crush-verify-${TIMESTAMP}"
BINARY="/tmp/crush-verify"

echo "=== LAUNCH ==="
echo "Building Crush from workspace..."
go build -o "${BINARY}" .

echo "Setting up test directory: ${TEST_DIR}"
mkdir -p "${TEST_DIR}/evidence"

echo "Verifying binary..."
"${BINARY}" --version

echo "✓ Launch complete"
echo "Binary: ${BINARY}"
echo "Test dir: ${TEST_DIR}"
echo "${TEST_DIR}" > /tmp/crush-verify-last-dir
