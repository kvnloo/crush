---
name: verify-crush
description: End-to-end verification harness for Crush. Use when the user asks to verify Crush functionality, test features, or run the verification suite. Performs Launch, Doctor, Drive, Evidence, and Cleanup phases for mapped features.
user-invocable: true
---

# Crush Verification Skill

This skill provides a comprehensive verification harness for Crush, the terminal-based AI coding assistant. It enables systematic testing of Crush features through a structured approach: Launch → Doctor → Drive → Evidence → Cleanup.

## When to Use

- User asks to verify Crush functionality
- Testing specific Crush features end-to-end
- Validating Crush after changes
- Running the verification suite
- Debugging feature regressions

## Architecture

The verification harness follows this structure:

```
.cursor/skills/verify-crush/
├── SKILL.md                    # This file (control orchestration)
├── features/
│   ├── README.md              # Feature map overview
│   ├── cli-basics.md          # CLI flags, help, version
│   ├── headless-run.md        # Non-interactive mode
│   ├── config-system.md       # Configuration loading
│   ├── session-management.md  # Session CRUD operations
│   └── tool-execution.md      # Built-in tools
└── helpers/
    ├── launch.sh              # Build/setup Crush binary
    ├── doctor.sh              # Pre-flight checks
    └── cleanup.sh             # Teardown verification state
```

## Phase Flow

### 1. Launch

**Purpose**: Prepare the Crush environment for testing.

**Actions**:
- Build Crush from source: `go build -o /tmp/crush-verify .`
- Set up isolated test directory: `/tmp/crush-verify-{timestamp}`
- Configure minimal test environment (no API keys required for basic tests)
- Verify binary is executable and responds to `--version`

**Script**: `helpers/launch.sh`

**Success Criteria**:
- Binary exists at `/tmp/crush-verify`
- `./crush-verify --version` returns exit code 0
- Test directory created and writable

### 2. Doctor

**Purpose**: Pre-flight health checks before driving features.

**Actions**:
- Verify Go toolchain availability: `go version`
- Check binary exists and is executable
- Validate test directory structure
- Confirm no conflicting processes
- Check for required dependencies (none needed for basic tests)

**Script**: `helpers/doctor.sh`

**Success Criteria**:
- All health checks pass
- Environment ready for feature testing

### 3. Drive

**Purpose**: Execute mapped features end-to-end.

**Actions**:
- For each feature in `features/`:
  - Parse the feature's "Driving it with <harness>" section
  - Execute test commands against `/tmp/crush-verify`
  - Capture stdout, stderr, and exit codes
  - Compare against expected behaviors
  - Save outputs to evidence directory

**Feature Selection**:
- Run all features: `crush verify` (default)
- Run specific feature: `crush verify --feature cli-basics`

**Test Execution**:
- Each feature defines its own test commands
- Tests should be deterministic and API-key-agnostic where possible
- For features requiring API keys, provide mock/stub paths

### 4. Evidence

**Purpose**: Capture verification results for review.

**Actions**:
- Save all command outputs to `/tmp/crush-verify-{timestamp}/evidence/`
- Structure: `evidence/{feature-name}/{test-name}.{out,err,exit}`
- Generate summary report: `evidence/SUMMARY.md`
- Include timestamps, exit codes, and pass/fail status

**Evidence Structure**:
```
evidence/
├── SUMMARY.md                  # Overall results
├── cli-basics/
│   ├── version.out
│   ├── version.exit
│   ├── help.out
│   └── help.exit
├── headless-run/
│   ├── simple-prompt.out
│   └── simple-prompt.exit
└── ...
```

**Preservation**:
- Evidence directory survives cleanup
- Path: `/tmp/crush-verify-{timestamp}/evidence/`
- Referenced in PR description with example outputs

### 5. Cleanup

**Purpose**: Tear down test state, preserving evidence.

**Actions**:
- Remove test binary: `/tmp/crush-verify`
- Remove test workspace: `/tmp/crush-verify-{timestamp}/*` (except `evidence/`)
- Restore original working directory
- Report evidence location

**Script**: `helpers/cleanup.sh`

**Preserved**:
- Evidence directory at `/tmp/crush-verify-{timestamp}/evidence/`
- SUMMARY.md with full results

## Feature Mapping

Features are documented in `features/*.md`. Each feature file includes:

1. **Sub-features**: What aspects are covered
2. **How to get to it**: User-facing description
3. **Driving it with <harness>**: Test commands and expected outcomes
4. **Gotchas**: Edge cases, limitations, or known issues

See `features/README.md` for the complete feature map.

## Usage Examples

### Run Full Verification

```bash
# From workspace root
crush verify
```

This executes all phases (Launch → Doctor → Drive → Evidence → Cleanup) for all mapped features.

### Verify Specific Feature

```bash
crush verify --feature cli-basics
```

### Re-run After Cleanup

```bash
# Evidence is preserved
cat /tmp/crush-verify-{timestamp}/evidence/SUMMARY.md
```

## Helpers

### Launch (`helpers/launch.sh`)

```bash
#!/usr/bin/env bash
set -euo pipefail

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
TEST_DIR="/tmp/crush-verify-${TIMESTAMP}"
BINARY="/tmp/crush-verify"

echo "=== LAUNCH ==="
echo "Building Crush..."
go build -o "${BINARY}" .

echo "Setting up test directory: ${TEST_DIR}"
mkdir -p "${TEST_DIR}/evidence"

echo "Verifying binary..."
"${BINARY}" --version

echo "✓ Launch complete"
echo "Binary: ${BINARY}"
echo "Test dir: ${TEST_DIR}"
```

### Doctor (`helpers/doctor.sh`)

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "=== DOCTOR ==="

echo "Checking Go toolchain..."
go version || { echo "✗ Go not found"; exit 1; }

echo "Checking binary..."
[[ -x /tmp/crush-verify ]] || { echo "✗ Binary not executable"; exit 1; }

echo "Checking test directory..."
[[ -d /tmp/crush-verify-* ]] || { echo "✗ Test directory missing"; exit 1; }

echo "✓ All checks passed"
```

### Cleanup (`helpers/cleanup.sh`)

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "=== CLEANUP ==="

echo "Removing binary..."
rm -f /tmp/crush-verify

echo "Finding test directory..."
TEST_DIR=$(ls -dt /tmp/crush-verify-* 2>/dev/null | head -1)

if [[ -n "${TEST_DIR}" ]]; then
    echo "Preserving evidence: ${TEST_DIR}/evidence/"
    echo "Removing test workspace (keeping evidence)..."
    find "${TEST_DIR}" -mindepth 1 -maxdepth 1 ! -name evidence -exec rm -rf {} +
    echo "✓ Cleanup complete"
    echo "Evidence: ${TEST_DIR}/evidence/"
else
    echo "No test directory found"
fi
```

## Implementation Notes

- **No API keys required**: Basic verification tests use `--help`, `--version`, and unit tests
- **Isolated execution**: All tests run in `/tmp/crush-verify-*` to avoid polluting workspace
- **Deterministic**: Tests should produce consistent results across runs
- **Evidence-driven**: Every test captures outputs for post-mortem analysis
- **CI-ready**: Can be integrated into CI pipelines (future work)

## Adding New Features

To add a new feature to the verification map:

1. Create `features/{feature-name}.md` following the template
2. Define sub-features, how to access, driving commands, and gotchas
3. Update `features/README.md` with the new feature entry
4. Test the feature with `crush verify --feature {feature-name}`

## Constraints

- Only operates on kvnloo/crush fork (never charmbracelet/crush upstream)
- Designed for headless/CLI verification (TUI testing out of scope)
- Evidence must survive cleanup for PR inclusion
- Tests prefer no-key paths (--help, unit tests) over live API calls

## Success Criteria

A verification run is successful when:

1. Launch: Binary built and smoke-tested
2. Doctor: All pre-flight checks pass
3. Drive: At least one feature tested end-to-end
4. Evidence: Outputs captured under named path
5. Cleanup: Test state removed, evidence preserved

## PR Integration

When verification is complete:

- Include evidence location in PR body
- Show example commands and exit codes
- Reference `features/README.md` for feature coverage
- Mention PER-1272 (Linear ticket)
- Note this is control skill + feature map contribution
