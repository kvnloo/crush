#!/usr/bin/env bash
set -euo pipefail

echo "=== DOCTOR ==="

echo "Checking Go toolchain..."
if ! go version; then
    echo "✗ Go not found"
    exit 1
fi
echo "✓ Go available: $(go version)"

echo "Checking binary..."
if [[ ! -x /tmp/crush-verify ]]; then
    echo "✗ Binary not executable at /tmp/crush-verify"
    exit 1
fi
echo "✓ Binary executable"

echo "Checking test directory..."
TEST_DIR=$(cat /tmp/crush-verify-last-dir 2>/dev/null || echo "")
if [[ -z "${TEST_DIR}" ]] || [[ ! -d "${TEST_DIR}" ]]; then
    echo "✗ Test directory missing"
    exit 1
fi
echo "✓ Test directory exists: ${TEST_DIR}"

echo "Checking evidence directory..."
if [[ ! -d "${TEST_DIR}/evidence" ]]; then
    echo "✗ Evidence directory missing"
    exit 1
fi
echo "✓ Evidence directory ready"

echo ""
echo "✓ All checks passed"
echo "Ready for feature testing"
