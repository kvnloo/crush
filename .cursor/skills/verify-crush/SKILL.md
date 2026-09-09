# Crush Verification Skill

This skill provides verification and diagnostic tools for the Crush codebase.

## Commands

### `crush doctor`

Performs comprehensive health checks across the Crush codebase:

- **Config System**: Validates crushrc parsing, provider resolution, and JSON fallback
- **Agent System**: Checks tool registration, prompt templates, and coordinator setup
- **Database**: Verifies SQLite schema, migrations, and sqlc generation
- **LSP Integration**: Tests LSP discovery, client lifecycle, and auto-start
- **MCP Integration**: Validates channel registration and client connections
- **UI/TUI**: Checks Bubble Tea model setup, styling system, and component structure
- **Hooks System**: Validates hook runner, decision types, and tool wrapping
- **Skills System**: Tests discovery, loading, and manager state
- **Testing**: Verifies test helpers, mock providers, and golden files

#### Usage

```bash
crush doctor
```

#### Output

```
🩺 Crush Doctor Report
═════════════════════════════════════════════════

✓ Config System        PASS
✓ Agent System         PASS
✗ Database             FAIL: migration 003 missing index
✓ LSP Integration      PASS
✓ MCP Integration      PASS
✓ UI/TUI               PASS
✓ Hooks System         PASS
✓ Skills System        PASS
✓ Testing              PASS

Overall: 8/9 checks passed
```

### `crush verify [feature]`

Runs deep verification on specific features:

- `config` - Config loading and provider resolution
- `agent` - Agent creation, tools, and execution
- `db` - Database schema and queries
- `lsp` - LSP client management
- `mcp` - MCP server integration
- `ui` - TUI components and styling
- `hooks` - Hook execution and decisions
- `skills` - Skill discovery and loading
- `all` - Run all feature verifications

#### Usage

```bash
crush verify config
crush verify all
```

## Feature Map

The verification system maps to these TUI agent surface areas:

### 1. **Editor/Input Surface** (`internal/ui/components/editor`)
- Command input and parsing
- Bash execution integration
- Permission prompts
- Skill invocation UI

### 2. **Message Display Surface** (`internal/ui/components/messages`)
- Agent response rendering
- Tool use visualization
- Markdown formatting
- Code syntax highlighting

### 3. **Status Surface** (`internal/ui/components/status`)
- Agent state display
- Provider/model information
- Session metadata
- Active tool indicators

### 4. **Session Management Surface** (`internal/ui/components/session`)
- Session list and selection
- Resume functionality
- Session history navigation

### 5. **Settings Surface** (`internal/ui/components/settings`)
- Provider configuration
- Permission management
- Theme selection
- Skill enable/disable

### 6. **Error/Warning Surface** (`internal/ui/components/feedback`)
- Error messages
- Warning banners
- Diagnostic information
- Recovery suggestions

## Implementation Notes

The `crush doctor` command is implemented in `internal/cmd/doctor.go` and uses
verification modules from `internal/verify/` to check each subsystem.

Each verification module follows this pattern:

```go
package verify

type Result struct {
    Name   string
    Status Status // PASS, FAIL, WARN, SKIP
    Error  error
    Details string
}

func VerifyConfig(ctx context.Context) Result
func VerifyAgent(ctx context.Context) Result
// ... etc
```

The doctor command aggregates results and presents them in a unified report.
