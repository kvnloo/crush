# Crush Feature Map

This directory contains mapped features for systematic verification of Crush. Each feature file documents what the feature does, how users access it, how to test it, and known gotchas.

## Feature Files

| Feature | File | Description |
|---------|------|-------------|
| CLI Basics | `cli-basics.md` | Version, help, dirs commands |
| Headless Run | `headless-run.md` | Non-interactive `crush run` mode |
| Config System | `config-system.md` | Configuration loading and precedence |
| Session Management | `session-management.md` | Session creation, listing, deletion |
| Tool Execution | `tool-execution.md` | Built-in tools (view, edit, bash) |

## Feature File Template

Each feature file follows this structure:

```markdown
# Feature Name

## Sub-features

List of specific capabilities covered by this feature.

## How to Get to It (User POV)

User-facing description of how to access and use this feature.

## Driving It with <harness>

Test commands, expected outputs, and success criteria for verification.

## Gotchas

Edge cases, limitations, known issues, or things that might trip up testing.
```

## Adding New Features

To add a new feature:

1. Create `features/{feature-name}.md` following the template
2. Fill in all four sections (Sub-features, How to get to it, Driving, Gotchas)
3. Update this README with an entry in the table above
4. Test the feature via the skill harness: `bash .cursor/skills/verify-crush/prove-one-feature.sh` (cli-basics) or follow Driving commands in the feature file after `helpers/launch.sh` + `helpers/doctor.sh`

## Coverage Status

- ✅ CLI Basics: Version, help, dirs
- ✅ Headless Run: Non-interactive mode
- ✅ Config System: Configuration loading
- ✅ Session Management: CRUD operations
- ✅ Tool Execution: Built-in tools

## Future Coverage

Potential features for future verification:

- Interactive TUI mode (requires screen recording)
- LSP integration (requires LSP server setup)
- MCP integration (requires MCP server)
- Provider configuration (requires API keys)
- Model selection and switching
- Permission system (allow/deny lists)
- Hook execution (pre-tool-use hooks)
- Skill system (loading and invoking skills)
- Git integration (commits, PRs)
- Multi-session workflows

## Testing Philosophy

- **API-key agnostic**: Prefer tests that don't require API keys
- **Deterministic**: Tests should produce consistent results
- **Isolated**: Each feature test is independent
- **Evidence-driven**: Capture all outputs for review
- **User-facing**: Test from the user's perspective, not internals
