# Config System

## Sub-features

- Configuration file discovery (precedence: ./.crushrc, ./crushrc, ~/.config/crush/crushrc)
- JSON config support (legacy crush.json)
- Bash-based crushrc format (using builtins: provider, model, lsp, mcp, options, permissions, hook)
- Environment variable expansion in config
- Schema validation
- Config reloading

## How to Get to It (User POV)

Users configure Crush via `crushrc` (Bash format) or `crush.json` (legacy):

**crushrc (preferred)**:
```bash
# ~/.config/crush/crushrc
provider add ollama --type ollama --base-url "http://localhost:11434/v1"
model add ollama/llama3.3 --name "Llama 3.3" --context-window 128000
permissions allow view edit
option debug true
```

**crush.json (legacy)**:
```json
{
  "$schema": "https://charm.land/crush.json",
  "providers": {
    "openai": {
      "type": "openai",
      "api_key": "$OPENAI_API_KEY"
    }
  }
}
```

Configuration is discovered in this order:
1. `./.crushrc` (project-specific)
2. `./crushrc`
3. `~/.config/crush/crushrc` (user-global)

## Driving It with <harness>

### Test: No Config (Default Behavior)

**Setup**:
```bash
# Ensure no config files exist in test environment
export HOME=${TEST_DIR}/fake-home
mkdir -p ${TEST_DIR}/fake-home/.config/crush
```

**Command**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Shows default config paths
- Exit code: 0

**Evidence**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs > ${TEST_DIR}/evidence/config-system/no-config.out 2>&1
echo $? > ${TEST_DIR}/evidence/config-system/no-config.exit
```

### Test: Empty crushrc Loads

**Setup**:
```bash
mkdir -p ${TEST_DIR}/fake-home/.config/crush
touch ${TEST_DIR}/fake-home/.config/crush/crushrc
```

**Command**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Loads without error
- Exit code: 0

**Evidence**:
```bash
touch ${TEST_DIR}/fake-home/.config/crush/crushrc
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs > ${TEST_DIR}/evidence/config-system/empty-crushrc.out 2>&1
echo $? > ${TEST_DIR}/evidence/config-system/empty-crushrc.exit
```

### Test: Simple crushrc with Options

**Setup**:
```bash
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
option debug true
option notifications disabled
EOF
```

**Command**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Parses config successfully
- Exit code: 0

**Evidence**:
```bash
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
option debug true
option notifications disabled
EOF
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs > ${TEST_DIR}/evidence/config-system/simple-crushrc.out 2>&1
echo $? > ${TEST_DIR}/evidence/config-system/simple-crushrc.exit
```

### Test: Config Precedence (Local Override)

**Setup**:
```bash
# Global config
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
option debug false
EOF

# Local config (should override)
cat > ${TEST_DIR}/test-workspace/crushrc <<'EOF'
option debug true
EOF
```

**Command**:
```bash
cd ${TEST_DIR}/test-workspace && HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Local config takes precedence
- Exit code: 0

**Evidence**:
```bash
mkdir -p ${TEST_DIR}/test-workspace
cat > ${TEST_DIR}/fake-home/.config/crush/crushrc <<'EOF'
option debug false
EOF
cat > ${TEST_DIR}/test-workspace/crushrc <<'EOF'
option debug true
EOF
(cd ${TEST_DIR}/test-workspace && HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs) > ${TEST_DIR}/evidence/config-system/precedence.out 2>&1
echo $? > ${TEST_DIR}/evidence/config-system/precedence.exit
```

### Test: JSON Config Still Works (Legacy)

**Setup**:
```bash
cat > ${TEST_DIR}/fake-home/.config/crush/crush.json <<'EOF'
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "debug": true
  }
}
EOF
```

**Command**:
```bash
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs
```

**Expected Output**:
- Loads JSON config
- Exit code: 0

**Evidence**:
```bash
cat > ${TEST_DIR}/fake-home/.config/crush/crush.json <<'EOF'
{
  "$schema": "https://charm.land/crush.json",
  "options": {
    "debug": true
  }
}
EOF
HOME=${TEST_DIR}/fake-home /tmp/crush-verify dirs > ${TEST_DIR}/evidence/config-system/json-config.out 2>&1
echo $? > ${TEST_DIR}/evidence/config-system/json-config.exit
```

## Gotchas

- **Bash execution**: crushrc is executed as Bash; syntax errors will fail config loading
- **Builtin scope**: Config builtins (provider, model, etc.) only work during config load, not in normal bash tool execution
- **Environment isolation**: Tests must set HOME to isolated directory to avoid loading user's real config
- **Deep merge**: Multiple config files are deep-merged; last occurrence wins for conflicting keys
- **Schema validation**: Invalid JSON configs may fail validation; crushrc has no schema enforcement beyond bash syntax
- **API key expansion**: $VAR and $(command) are expanded at config load time
- **Security**: crushrc runs arbitrary bash; never load untrusted configs
- **XDG variables**: Crush respects XDG_CONFIG_HOME; tests should set it explicitly if needed
