package mcp

import (
	"io"
	"testing"

	"github.com/charmbracelet/crush/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func TestCreateTransport_StdioDiscardsChildStderr(t *testing.T) {
	t.Parallel()

	m := config.MCPConfig{Type: config.MCPStdio, Command: "echo", Args: []string{"hi"}}
	tr, _, err := createTransport(t.Context(), nil, "test", m, shellResolverWithPath(t, nil))
	require.NoError(t, err)

	ct, ok := tr.(*mcp.CommandTransport)
	require.True(t, ok, "expected CommandTransport, got %T", tr)
	require.Equal(t, io.Discard, ct.Command.Stderr)
}
