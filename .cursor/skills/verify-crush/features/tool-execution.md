# Tool Execution

## Sub-features

- Built-in tools: view, edit, write, bash, grep, glob, diagnostics, lsp_*, job_*, crush_info
- Tool permission system (allow/deny lists)
- Tool execution without API calls (unit tests validate tool logic)
- Tool self-documentation (each tool has .md description file)

## How to Get to It (User POV)

Crush provides tools that the LLM agent can invoke during a session. Users don't directly call these tools; instead, they're invoked by the agent based on the conversation context.

However, users can configure which tools are allowed:

```bash
# In crushrc
permissions allow view edit grep
permissions deny bash sourcegraph
```

Tools are implemented in `internal/agent/tools/` and documented in corresponding `.md` files.

## Driving It with <harness>

Since tools are invoked by the agent (which requires an LLM), we test tool logic via unit tests rather than end-to-end execution.

### Test: Tool Unit Tests Exist

**Command**:
```bash
go test ./internal/agent/tools -list . | grep -i test | wc -l
```

**Expected Output**:
- Shows count of test cases
- Should be > 0
- Exit code: 0

**Evidence**:
```bash
go test ./internal/agent/tools -list . | grep -i test > ${TEST_DIR}/evidence/tool-execution/test-list.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/test-list.exit
```

### Test: Run Tool Unit Tests

**Command**:
```bash
go test ./internal/agent/tools -run TestCrushInfo_MinimalConfig -v
```

**Expected Output**:
- Test passes (PASS)
- Exit code: 0

**Evidence**:
```bash
go test ./internal/agent/tools -run TestCrushInfo_MinimalConfig -v > ${TEST_DIR}/evidence/tool-execution/crush-info-test.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/crush-info-test.exit
```

### Test: Run Bash Tool Tests

**Command**:
```bash
go test ./internal/agent/tools -run TestBashTool -v
```

**Expected Output**:
- Tests pass (PASS)
- Exit code: 0

**Evidence**:
```bash
go test ./internal/agent/tools -run TestBashTool -v > ${TEST_DIR}/evidence/tool-execution/bash-tool-tests.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/bash-tool-tests.exit
```

### Test: Tool Documentation Exists

**Command**:
```bash
find internal/agent/tools -name "*.md" -type f | wc -l
```

**Expected Output**:
- Shows count of tool documentation files
- Should be > 10 (one per tool)
- Exit code: 0

**Evidence**:
```bash
find internal/agent/tools -name "*.md" -type f > ${TEST_DIR}/evidence/tool-execution/tool-docs-list.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/tool-docs-list.exit
```

### Test: Tool Permissions Config Parsing

**Setup**:
```bash
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
permissions allow view edit
permissions deny bash
EOF
```

**Command**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Config loads without error
- Exit code: 0

**Evidence**:
```bash
mkdir -p ${TEST_DIR}/fake-home/.config/crush
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
permissions allow view edit
permissions deny bash
EOF
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs > ${TEST_DIR}/evidence/tool-execution/permissions-config.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/permissions-config.exit
```

### Test: YOLO Mode Flag

**Command**:
```bash
/tmp/crush-verify --yolo --help
```

**Expected Output**:
- Accepts --yolo flag
- Help displays normally
- Exit code: 0

**Evidence**:
```bash
/tmp/crush-verify --yolo --help > ${TEST_DIR}/evidence/tool-execution/yolo-flag.out 2>&1
echo $? > ${TEST_DIR}/evidence/tool-execution/yolo-flag.exit
```

## Gotchas

- **Agent-invoked**: Tools are not directly callable by users; they're invoked by the LLM agent during a session
- **Permission prompts**: Without --yolo, Crush prompts for permission before executing tools
- **Unit tests**: Tool logic is tested via unit tests in `internal/agent/tools/`; we rely on these rather than end-to-end LLM execution
- **No live execution**: Without API keys, we can't test tools in a live agent session
- **Tool context**: Some tools require context (session, message, model) that's only available during agent execution
- **Bash tool**: The bash tool executes shell commands; it has special permission handling and auto-backgrounding logic
- **MCP tools**: MCP servers expose additional tools dynamically; testing those requires MCP server setup
- **LSP tools**: LSP tools require an LSP server running; not tested in basic verification
- **Job management**: job_output and job_kill manage background jobs; they depend on bash tool execution
- **Documentation format**: Tool .md files use a specific format for LLM consumption; they're not user-facing docs
