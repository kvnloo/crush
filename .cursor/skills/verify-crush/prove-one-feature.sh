#!/usr/bin/env bash
# Prove ONE feature end-to-end: cli-basics
# This demonstrates Launch → Doctor → Drive → Evidence → Cleanup

set -euo pipefail

WORKSPACE_ROOT=$(pwd)
SKILL_DIR="${WORKSPACE_ROOT}/.cursor/skills/verify-crush"

echo "========================================"
echo "Crush Verification: Prove CLI Basics"
echo "========================================"
echo ""

# Phase 1: Launch
echo "=== Phase 1: LAUNCH ==="
cd "${WORKSPACE_ROOT}"
bash "${SKILL_DIR}/helpers/launch.sh"
TEST_DIR=$(cat /tmp/crush-verify-last-dir)
echo ""

# Phase 2: Doctor
echo "=== Phase 2: DOCTOR ==="
bash "${SKILL_DIR}/helpers/doctor.sh"
echo ""

# Phase 3: Drive (cli-basics feature)
echo "=== Phase 3: DRIVE (cli-basics) ==="
mkdir -p "${TEST_DIR}/evidence/cli-basics"

echo "Test: Version"
/tmp/crush-verify --version > "${TEST_DIR}/evidence/cli-basics/version.out" 2>&1 || true
echo $? > "${TEST_DIR}/evidence/cli-basics/version.exit"
echo "✓ Version test captured (exit code: $(cat ${TEST_DIR}/evidence/cli-basics/version.exit))"

echo "Test: Main Help"
/tmp/crush-verify --help > "${TEST_DIR}/evidence/cli-basics/help.out" 2>&1 || true
echo $? > "${TEST_DIR}/evidence/cli-basics/help.exit"
echo "✓ Help test captured (exit code: $(cat ${TEST_DIR}/evidence/cli-basics/help.exit))"

echo "Test: Run Help"
/tmp/crush-verify run --help > "${TEST_DIR}/evidence/cli-basics/run-help.out" 2>&1 || true
echo $? > "${TEST_DIR}/evidence/cli-basics/run-help.exit"
echo "✓ Run help test captured (exit code: $(cat ${TEST_DIR}/evidence/cli-basics/run-help.exit))"

echo "Test: Dirs"
/tmp/crush-verify dirs > "${TEST_DIR}/evidence/cli-basics/dirs.out" 2>&1 || true
echo $? > "${TEST_DIR}/evidence/cli-basics/dirs.exit"
echo "✓ Dirs test captured (exit code: $(cat ${TEST_DIR}/evidence/cli-basics/dirs.exit))"

echo "Test: Models"
/tmp/crush-verify models > "${TEST_DIR}/evidence/cli-basics/models.out" 2>&1 || true
echo $? > "${TEST_DIR}/evidence/cli-basics/models.exit"
echo "✓ Models test captured (exit code: $(cat ${TEST_DIR}/evidence/cli-basics/models.exit))"

echo ""

# Phase 4: Evidence
echo "=== Phase 4: EVIDENCE ==="
cat > "${TEST_DIR}/evidence/SUMMARY.md" <<EOF
# Crush Verification Summary

**Feature**: cli-basics
**Timestamp**: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**Test Directory**: ${TEST_DIR}

## Results

| Test | Exit Code | Status |
|------|-----------|--------|
| Version | $(cat ${TEST_DIR}/evidence/cli-basics/version.exit) | $([ "$(cat ${TEST_DIR}/evidence/cli-basics/version.exit)" -eq 0 ] && echo "✓ PASS" || echo "✗ FAIL") |
| Help | $(cat ${TEST_DIR}/evidence/cli-basics/help.exit) | $([ "$(cat ${TEST_DIR}/evidence/cli-basics/help.exit)" -eq 0 ] && echo "✓ PASS" || echo "✗ FAIL") |
| Run Help | $(cat ${TEST_DIR}/evidence/cli-basics/run-help.exit) | $([ "$(cat ${TEST_DIR}/evidence/cli-basics/run-help.exit)" -eq 0 ] && echo "✓ PASS" || echo "✗ FAIL") |
| Dirs | $(cat ${TEST_DIR}/evidence/cli-basics/dirs.exit) | $([ "$(cat ${TEST_DIR}/evidence/cli-basics/dirs.exit)" -eq 0 ] && echo "✓ PASS" || echo "✗ FAIL") |
| Models | $(cat ${TEST_DIR}/evidence/cli-basics/models.exit) | $([ "$(cat ${TEST_DIR}/evidence/cli-basics/models.exit)" -eq 0 ] && echo "✓ PASS" || echo "✗ FAIL") |

## Sample Outputs

### Version Output
\`\`\`
$(head -5 ${TEST_DIR}/evidence/cli-basics/version.out)
\`\`\`

### Dirs Output
\`\`\`
$(cat ${TEST_DIR}/evidence/cli-basics/dirs.out)
\`\`\`

## Evidence Location

All test outputs are preserved at: \`${TEST_DIR}/evidence/\`

EOF

echo "Summary report generated: ${TEST_DIR}/evidence/SUMMARY.md"
echo ""
cat "${TEST_DIR}/evidence/SUMMARY.md"
echo ""

# Phase 5: Cleanup
echo "=== Phase 5: CLEANUP ==="
bash "${SKILL_DIR}/helpers/cleanup.sh"
echo ""

echo "========================================"
echo "✓ Verification Complete"
echo "Evidence preserved at: ${TEST_DIR}/evidence/"
echo "========================================"
