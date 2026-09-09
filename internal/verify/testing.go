package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifyTesting checks the testing infrastructure.
func VerifyTesting(ctx context.Context) Result {
	// Check for test files in key packages.
	testPaths := []string{
		"internal/config",
		"internal/agent",
		"internal/db",
	}

	foundTests := 0
	for _, pkg := range testPaths {
		if hasTestFiles(pkg) {
			foundTests++
		}
	}

	if foundTests == 0 {
		return Result{
			Name:    "Testing",
			Status:  StatusWarn,
			Details: "no test files found",
		}
	}

	return Result{
		Name:    "Testing",
		Status:  StatusPass,
		Details: "tests present",
	}
}

func hasTestFiles(dir string) bool {
	candidates := []string{dir, filepath.Join("..", dir), filepath.Join("../..", dir)}
	for _, c := range candidates {
		entries, err := os.ReadDir(c)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == ".go" {
				if len(entry.Name()) > 8 && entry.Name()[len(entry.Name())-8:] == "_test.go" {
					return true
				}
			}
		}
	}
	return false
}
