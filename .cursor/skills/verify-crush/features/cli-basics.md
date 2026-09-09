# CLI Basics

## Sub-features

- `--version`: Display Crush version information
- `--help`: Display help text for commands
- `dirs`: Show config and data directories
- `models`: List available models from providers (no API keys needed)

## How to Get to It (User POV)

Users interact with these basic CLI commands from the terminal:

```bash
# Check version
crush --version

# Get help on main command
crush --help

# Get help on specific subcommand
crush run --help

# Show directories
crush dirs

# List models (shows configured and discovered providers)
crush models
```

These are the first things users encounter when exploring Crush. They provide essential information about the installation, available commands, and configuration locations.

## Driving It with <harness>

### Test: Version Display

**Command**:
```bash
/tmp/crush-verify --version
```

**Expected Output**:
- Stdout contains "crush version"
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify --version > ${TEST_DIR}/evidence/cli-basics/version.out 2>&1
echo $? > ${TEST_DIR}/evidence/cli-basics/version.exit
```

### Test: Main Help

**Command**:
```bash
/tmp/crush-verify --help
```

**Expected Output**:
- Stdout contains "USAGE", "COMMANDS", "FLAGS"
- Mentions key commands: run, models, session
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify --help > ${TEST_DIR}/evidence/cli-basics/help.out 2>&1
echo $? > ${TEST_DIR}/evidence/cli-basics/help.exit
```

### Test: Subcommand Help

**Command**:
```bash
/tmp/crush-verify run --help
```

**Expected Output**:
- Stdout contains "Run a single prompt"
- Shows flags: --model, --session, --quiet
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify run --help > ${TEST_DIR}/evidence/cli-basics/run-help.out 2>&1
echo $? > ${TEST_DIR}/evidence/cli-basics/run-help.exit
```

### Test: Dirs Command

**Command**:
```bash
/tmp/crush-verify dirs
```

**Expected Output**:
- Stdout contains paths (e.g., ~/.config/crush, ~/.local/share/crush)
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify dirs > ${TEST_DIR}/evidence/cli-basics/dirs.out 2>&1
echo $? > ${TEST_DIR}/evidence/cli-basics/dirs.exit
```

### Test: Models Command (No Config)

**Command**:
```bash
/tmp/crush-verify models
```

**Expected Output**:
- Stdout may show "No providers configured" or list built-in providers
- Exit code: 0 (should not fail even without config)

**Evidence**:
```bash
/tmp/crush-verify models > ${TEST_DIR}/evidence/cli-basics/models.out 2>&1
echo $? > ${TEST_DIR}/evidence/cli-basics/models.exit
```

## Gotchas

- **Version format**: The version string format may vary (e.g., "devel" during development vs. "v1.2.3" in releases)
- **Help text length**: Help output is verbose; expect multi-page output
- **Dirs platform-specific**: Paths differ on Linux/macOS/Windows (test accounts for this)
- **Models without config**: The `models` command should gracefully handle missing configuration and not crash
- **Stdout vs stderr**: Some help text may go to stderr; capture both streams
- **Terminal width**: Help formatting may vary with terminal width; harness uses default
