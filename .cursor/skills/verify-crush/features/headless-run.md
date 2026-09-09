# Headless Run

## Sub-features

- `crush run` with prompt arguments
- `crush run` with stdin piping
- `--quiet` flag to suppress spinner
- `--verbose` flag to show logs
- Exit code handling (0 for success, non-zero for failures)

## How to Get to It (User POV)

The `crush run` command allows non-interactive, headless operation of Crush:

```bash
# Run with arguments
crush run "What is 2+2?"

# Pipe input
echo "Summarize this" | crush run "Process stdin"

# Quiet mode (no spinner)
crush run --quiet "Generate a README"

# Verbose mode (show logs)
crush run --verbose "Explain this code"

# Continue a session
crush run --session abc123 "Follow up question"
```

This is critical for scripting, CI/CD integration, and non-interactive workflows.

## Driving It with <harness>

### Test: Simple Headless Run (No API Key)

Since we don't have API keys configured, we test the **command parsing and error handling** rather than actual LLM execution.

**Command**:
```bash
/tmp/crush-verify run "test prompt" 2>&1 || true
```

**Expected Behavior**:
- Command should parse successfully
- Should fail gracefully with "no provider configured" or similar error
- Should NOT crash or hang
- Exit code: non-zero (expected failure due to missing config)

**Evidence**:
```bash
/tmp/crush-verify run "test prompt" > ${TEST_DIR}/evidence/headless-run/simple-run.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/headless-run/simple-run.exit
```

### Test: Help Flag Works

**Command**:
```bash
/tmp/crush-verify run --help
```

**Expected Output**:
- Stdout contains "Run a single prompt in non-interactive mode"
- Shows usage examples
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify run --help > ${TEST_DIR}/evidence/headless-run/help.out 2>&1
echo $? > ${TEST_DIR}/evidence/headless-run/help.exit
```

### Test: Quiet Flag Parsing

**Command**:
```bash
/tmp/crush-verify run --quiet "test" 2>&1 || true
```

**Expected Behavior**:
- Accepts --quiet flag without error
- Should fail due to missing config (not due to flag parsing)

**Evidence**:
```bash
/tmp/crush-verify run --quiet "test" > ${TEST_DIR}/evidence/headless-run/quiet-flag.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/headless-run/quiet-flag.exit
```

### Test: Verbose Flag Parsing

**Command**:
```bash
/tmp/crush-verify run --verbose "test" 2>&1 || true
```

**Expected Behavior**:
- Accepts --verbose flag without error
- May show additional debug output

**Evidence**:
```bash
/tmp/crush-verify run --verbose "test" > ${TEST_DIR}/evidence/headless-run/verbose-flag.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/headless-run/verbose-flag.exit
```

### Test: Stdin Pipe Handling

**Command**:
```bash
echo "input data" | /tmp/crush-verify run "process this" 2>&1 || true
```

**Expected Behavior**:
- Command accepts stdin without hanging
- Fails gracefully due to missing config

**Evidence**:
```bash
echo "input data" | /tmp/crush-verify run "process this" > ${TEST_DIR}/evidence/headless-run/stdin-pipe.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/headless-run/stdin-pipe.exit
```

## Gotchas

- **No API keys**: Without configured providers, `crush run` will fail after parsing. This is expected. We're testing that it parses correctly and fails gracefully, not that it can actually execute LLM requests.
- **Timeout behavior**: If the command hangs waiting for input, kill after 5 seconds
- **Spinner in non-TTY**: The spinner should be suppressed when stdout is not a TTY
- **Exit codes**: Exit code 0 means success; non-zero means failure. In our no-key tests, non-zero is expected.
- **Session flags**: `--session` and `--continue` require existing sessions; they'll fail gracefully if sessions don't exist
- **Model flag**: `--model` requires a configured provider; will fail without one
- **Verbose output**: May include debug logs, timestamps, and internal details
