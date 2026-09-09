package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyMCP checks MCP integration.
func VerifyMCP(ctx context.Context) Result {
	// Check that MCP package exists.
	mcpDir := findMCPDir()
	if mcpDir == "" {
		return Result{
			Name:   "MCP Integration",
			Status: StatusWarn,
			Details: "mcp directory not found",
		}
	}

	// Check for key files.
	requiredFiles := []string{"client.go"}
	for _, file := range requiredFiles {
		path := filepath.Join(mcpDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return Result{
				Name:    "MCP Integration",
				Status:  StatusWarn,
				Details: file + " not found",
			}
		}
	}

	return Result{
		Name:   "MCP Integration",
		Status: StatusPass,
	}
}

func findMCPDir() string {
	candidates := []string{
		"internal/mcp",
		"../mcp",
		"../../mcp",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
