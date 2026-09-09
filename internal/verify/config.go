package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyConfig checks the configuration system.
func VerifyConfig(ctx context.Context) Result {
	// Check that config package can be imported and has expected structure.
	configDir := findConfigDir()
	if configDir == "" {
		return Result{
			Name:    "Config System",
			Status:  StatusWarn,
			Details: "config directory not found",
		}
	}

	// Check for key files.
	requiredFiles := []string{"config.go", "load.go", "provider.go"}
	for _, file := range requiredFiles {
		path := filepath.Join(configDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return Result{
				Name:    "Config System",
				Status:  StatusWarn,
				Details: file + " not found",
			}
		}
	}

	// Check shell config support (it's in internal/shellconfig, not internal/config/shellconfig).
	shellConfigDir := findShellConfigDir()
	if shellConfigDir == "" {
		return Result{
			Name:    "Config System",
			Status:  StatusWarn,
			Details: "shellconfig directory not found",
		}
	}

	return Result{
		Name:   "Config System",
		Status: StatusPass,
	}
}

func findConfigDir() string {
	candidates := []string{
		"internal/config",
		"../config",
		"../../config",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}

func findShellConfigDir() string {
	candidates := []string{
		"internal/shellconfig",
		"../shellconfig",
		"../../shellconfig",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
