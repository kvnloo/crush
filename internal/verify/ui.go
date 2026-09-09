package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyUI checks the UI/TUI system.
func VerifyUI(ctx context.Context) Result {
	// Check that UI package exists.
	uiDir := findUIDir()
	if uiDir == "" {
		return Result{
			Name:   "UI/TUI",
			Status: StatusWarn,
			Details: "ui directory not found",
		}
	}

	// Check for key subdirectories.
	requiredDirs := []string{"components", "styles", "model"}
	for _, dir := range requiredDirs {
		path := filepath.Join(uiDir, dir)
		if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
			return Result{
				Name:    "UI/TUI",
				Status:  StatusWarn,
				Details: dir + " directory not found",
			}
		}
	}

	return Result{
		Name:   "UI/TUI",
		Status: StatusPass,
	}
}

func findUIDir() string {
	candidates := []string{
		"internal/ui",
		"../ui",
		"../../ui",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
