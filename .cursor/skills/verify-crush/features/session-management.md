# Session Management

## Sub-features

- Session creation (implicit on first run)
- Session listing (`crush session list`)
- Session deletion (`crush session delete`)
- Session continuation (`--session` flag)
- Session metadata (ID, timestamps, message count)
- Session storage (SQLite database)

## How to Get to It (User POV)

Users interact with sessions via the `crush session` subcommand:

```bash
# List all sessions
crush session list

# Delete a session
crush session delete <session-id>

# Delete all sessions
crush session delete --all

# Continue a specific session
crush run --session <session-id> "Follow up question"

# Continue most recent session
crush run --continue "Follow up"
```

Sessions are automatically created when running Crush interactively or via `crush run`. Each session stores the conversation history in SQLite.

## Driving It with <harness>

### Test: Session List (Empty)

**Setup**:
```bash
# Use isolated data directory
export CRUSH_DATA_DIR=${TEST_DIR}/data
mkdir -p ${TEST_DIR}/data
```

**Command**:
```bash
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session list
```

**Expected Output**:
- Shows no sessions (or empty list)
- Exit code: 0

**Evidence**:
```bash
mkdir -p ${TEST_DIR}/data
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session list > ${TEST_DIR}/evidence/session-management/empty-list.out 2>&1
echo $? > ${TEST_DIR}/evidence/session-management/empty-list.exit
```

### Test: Session List Help

**Command**:
```bash
/tmp/crush-verify session list --help
```

**Expected Output**:
- Shows usage information for listing sessions
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify session list --help > ${TEST_DIR}/evidence/session-management/list-help.out 2>&1
echo $? > ${TEST_DIR}/evidence/session-management/list-help.exit
```

### Test: Session Delete Help

**Command**:
```bash
/tmp/crush-verify session delete --help
```

**Expected Output**:
- Shows usage information for deleting sessions
- Mentions --all flag
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify session delete --help > ${TEST_DIR}/evidence/session-management/delete-help.out 2>&1
echo $? > ${TEST_DIR}/evidence/session-management/delete-help.exit
```

### Test: Session Subcommand Help

**Command**:
```bash
/tmp/crush-verify session --help
```

**Expected Output**:
- Shows available session subcommands (list, delete)
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify session --help > ${TEST_DIR}/evidence/session-management/session-help.out 2>&1
echo $? > ${TEST_DIR}/evidence/session-management/session-help.exit
```

### Test: Delete Non-Existent Session (Error Handling)

**Command**:
```bash
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session delete nonexistent-id 2>&1 || true
```

**Expected Output**:
- Should handle gracefully (error message or no-op)
- Exit code: non-zero (expected failure)

**Evidence**:
```bash
mkdir -p ${TEST_DIR}/data
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session delete nonexistent-id > ${TEST_DIR}/evidence/session-management/delete-nonexistent.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/session-management/delete-nonexistent.exit
```

### Test: Delete All Sessions (Empty)

**Command**:
```bash
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session delete --all 2>&1 || true
```

**Expected Output**:
- Should complete without error (even if no sessions exist)
- Exit code: 0

**Evidence**:
```bash
mkdir -p ${TEST_DIR}/data
CRUSH_DATA_DIR=${TEST_DIR}/data /tmp/crush-verify session delete --all > ${TEST_DIR}/evidence/session-management/delete-all-empty.out 2>&1 || true
echo $? > ${TEST_DIR}/evidence/session-management/delete-all-empty.exit
```

## Gotchas

- **Database location**: Sessions are stored in `~/.local/share/crush/crush.db` by default; use `--data-dir` or `CRUSH_DATA_DIR` to isolate for testing
- **Session IDs**: Session IDs are UUIDs; they're not sequential
- **Implicit creation**: Sessions are created automatically on first message; there's no explicit "create session" command
- **Session state**: Deleting a session removes all messages; this is irreversible
- **Concurrent access**: SQLite handles concurrent access, but race conditions are possible with multiple clients
- **No session rename**: Sessions cannot be renamed; they have IDs only
- **List pagination**: Large session lists may need pagination (not currently tested)
- **Session export**: There's no built-in session export feature (future work)
- **--continue flag**: The `--continue` flag finds the most recent session; it may behave unexpectedly if sessions are created/deleted concurrently
