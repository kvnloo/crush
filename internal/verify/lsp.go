package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyLSP checks LSP integration.
func VerifyLSP(ctx context.Context) Result {
	// Check that LSP package exists and has expected structure.
	lspDir := findLSPDir()
	if lspDir == "" {
		return Result{
			Name:   "LSP Integration",
			Status: StatusWarn,
			Error:  nil,
			Details: "lsp directory not found",
		}
	}

	// Check for key files.
	requiredFiles := []string{"manager.go", "client.go"}
	for _, file := range requiredFiles {
		path := filepath.Join(lspDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return Result{
				Name:    "LSP Integration",
				Status:  StatusWarn,
				Details: file + " not found",
			}
		}
	}

	return Result{
		Name:   "LSP Integration",
		Status: StatusPass,
	}
}

func findLSPDir() string {
	candidates := []string{
		"internal/lsp",
		"../lsp",
		"../../lsp",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
