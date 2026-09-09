package verify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// VerifyAgent checks the agent system.
func VerifyAgent(ctx context.Context) Result {
	// Check that tool files exist.
	toolsDir := findToolsDir()
	if toolsDir == "" {
		return Result{
			Name:   "Agent System",
			Status: StatusWarn,
			Details: "tools directory not found",
		}
	}

	// Count tool implementation files.
	entries, err := os.ReadDir(toolsDir)
	if err != nil {
		return Result{
			Name:   "Agent System",
			Status: StatusFail,
			Error:  fmt.Errorf("failed to read tools directory: %w", err),
		}
	}

	toolCount := 0
	for _, entry := range entries {
		name := entry.Name()
		// Count .go files that aren't tests.
		if !entry.IsDir() && filepath.Ext(name) == ".go" && len(name) > 8 && name[len(name)-8:] != "_test.go" {
			toolCount++
		}
	}

	if toolCount == 0 {
		return Result{
			Name:   "Agent System",
			Status: StatusFail,
			Error:  fmt.Errorf("no tool files found"),
		}
	}

	return Result{
		Name:    "Agent System",
		Status:  StatusPass,
		Details: fmt.Sprintf("%d tool files", toolCount),
	}
}

func findToolsDir() string {
	candidates := []string{
		"internal/agent/tools",
		"../agent/tools",
		"../../agent/tools",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
