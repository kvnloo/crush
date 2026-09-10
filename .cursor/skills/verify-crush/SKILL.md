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
├── prove-one-feature.sh        # Packaged cli-basics Launch→…→Cleanup
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
- Prove cli-basics end-to-end: `bash .cursor/skills/verify-crush/prove-one-feature.sh`
- Other features: run `helpers/launch.sh` → `helpers/doctor.sh`, then the feature file's Driving commands against `/tmp/crush-verify`
- There is no product `crush verify` subcommand; the harness lives under `.cursor/skills/verify-crush/`

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
# From workspace root (cli-basics proof today)
bash .cursor/skills/verify-crush/prove-one-feature.sh
```

This executes Launch → Doctor → Drive → Evidence → Cleanup for the cli-basics feature. Other features use the same helpers plus their feature-file Driving sections.

### Verify Specific Feature

```bash
# cli-basics packaged proof:
bash .cursor/skills/verify-crush/prove-one-feature.sh

# Or manually for any mapped feature:
bash .cursor/skills/verify-crush/helpers/launch.sh
bash .cursor/skills/verify-crush/helpers/doctor.sh
# then run the Driving commands from features/<name>.md
```

### Re-run After Cleanup

```bash
# Evidence is preserved
cat /tmp/crush-verify-{timestamp}/evidence/SUMMARY.md
```

## Helpers

Scripts under `helpers/` are the source of truth (do not copy stale snippets from this file):

- `helpers/launch.sh` — `go build -o /tmp/crush-verify .`, create `/tmp/crush-verify-$TIMESTAMP`, write path to `/tmp/crush-verify-last-dir`, smoke `--version`
- `helpers/doctor.sh` — require `go`, executable `/tmp/crush-verify`, and the last test/evidence dirs from `/tmp/crush-verify-last-dir`
- `helpers/cleanup.sh` — remove `/tmp/crush-verify` and workspace contents except `evidence/`
- `prove-one-feature.sh` — packaged Launch→Doctor→Drive→Evidence→Cleanup for **cli-basics**

```bash
bash .cursor/skills/verify-crush/helpers/launch.sh
bash .cursor/skills/verify-crush/helpers/doctor.sh
bash .cursor/skills/verify-crush/prove-one-feature.sh
bash .cursor/skills/verify-crush/helpers/cleanup.sh
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
4. Drive it with `helpers/launch.sh` + `helpers/doctor.sh` and the feature file's Driving section (cli-basics: `prove-one-feature.sh`)

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
