#!/usr/bin/env bash
set -euo pipefail

echo "=== CLEANUP ==="

echo "Removing binary..."
if [[ -f /tmp/crush-verify ]]; then
    rm -f /tmp/crush-verify
    echo "✓ Binary removed"
else
    echo "⊘ Binary not found (already removed?)"
fi

echo "Finding test directory..."
TEST_DIR=$(cat /tmp/crush-verify-last-dir 2>/dev/null || echo "")

if [[ -n "${TEST_DIR}" ]] && [[ -d "${TEST_DIR}" ]]; then
    echo "Preserving evidence: ${TEST_DIR}/evidence/"
    echo "Removing test workspace (keeping evidence)..."
    
    # Remove everything except evidence directory
    find "${TEST_DIR}" -mindepth 1 -maxdepth 1 ! -name evidence -exec rm -rf {} + 2>/dev/null || true
    
    echo "✓ Cleanup complete"
    echo "Evidence preserved at: ${TEST_DIR}/evidence/"
    
    # Clean up tracking file
    rm -f /tmp/crush-verify-last-dir
else
    echo "⊘ No test directory found to clean up"
fi

echo ""
echo "Cleanup finished"
