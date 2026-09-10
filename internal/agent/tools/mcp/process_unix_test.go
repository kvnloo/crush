//go:build !windows

package mcp

import (
	"os"
	"testing"

	"github.com/charmbracelet/crush/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

// TestCreateTransport_StdioProcessGroup pins that a stdio MCP child is spawned
// as its own process-group leader with a cancel hook wired up. This is what
// lets Crush reap a server's descendant processes (e.g. signal-cli launched by
// signal-mcp) when the session context is cancelled, instead of orphaning them.
func TestCreateTransport_StdioProcessGroup(t *testing.T) {
	t.Parallel()

	m := config.MCPConfig{Type: config.MCPStdio, Command: "echo", Args: []string{"hi"}}
	tr, _, err := createTransport(t.Context(), nil, "test", m, shellResolverWithPath(t, nil))
	require.NoError(t, err)

	ct, ok := tr.(*mcp.CommandTransport)
	require.True(t, ok, "expected CommandTransport, got %T", tr)
	require.NotNil(t, ct.Command.SysProcAttr, "stdio child must set SysProcAttr")
	require.True(t, ct.Command.SysProcAttr.Setpgid,
		"stdio child must lead its own process group so cancellation kills the whole tree")
	require.NotNil(t, ct.Command.Cancel,
		"stdio child must set a Cancel hook that kills the process group")
}

// TestCreateTransport_StdioStderrDiscard pins that a stdio MCP child's stderr
// is discarded rather than inherited from the parent. This prevents stderr
// output (e.g. MallocStackLogging warnings, debug logs) from clobbering the
// TUI, since the JSON-RPC protocol only uses stdout.
func TestCreateTransport_StdioStderrDiscard(t *testing.T) {
	t.Parallel()

	m := config.MCPConfig{Type: config.MCPStdio, Command: "echo", Args: []string{"hi"}}
	tr, _, err := createTransport(t.Context(), nil, "test", m, shellResolverWithPath(t, nil))
	require.NoError(t, err)

	ct, ok := tr.(*mcp.CommandTransport)
	require.True(t, ok, "expected CommandTransport, got %T", tr)
	require.NotNil(t, ct.Command.Stderr,
		"stdio child Stderr must be set (not nil) to prevent inheriting parent's stderr")
	require.NotEqual(t, ct.Command.Stderr, os.Stderr,
		"stdio child Stderr must not be os.Stderr to prevent TUI clobbering")
}
