package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyHooks checks the hooks system.
func VerifyHooks(ctx context.Context) Result {
	// Check that hooks package exists.
	hooksDir := findHooksDir()
	if hooksDir == "" {
		return Result{
			Name:   "Hooks System",
			Status: StatusWarn,
			Details: "hooks directory not found",
		}
	}

	// Check for key files.
	requiredFiles := []string{"hooks.go", "runner.go", "input.go"}
	for _, file := range requiredFiles {
		path := filepath.Join(hooksDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return Result{
				Name:    "Hooks System",
				Status:  StatusWarn,
				Details: file + " not found",
			}
		}
	}

	return Result{
		Name:   "Hooks System",
		Status: StatusPass,
	}
}

func findHooksDir() string {
	candidates := []string{
		"internal/hooks",
		"../hooks",
		"../../hooks",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
