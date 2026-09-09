package verify

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/crush/internal/db"
)

// VerifyDatabase checks the database layer.
func VerifyDatabase(ctx context.Context) Result {
	tempDir, err := os.MkdirTemp("", "crush-verify-db-*")
	if err != nil {
		return Result{
			Name:   "Database",
			Status: StatusFail,
			Error:  fmt.Errorf("failed to create temp dir: %w", err),
		}
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "crush.db")
	conn, err := db.Connect(ctx, tempDir)
	if err != nil {
		return Result{
			Name:   "Database",
			Status: StatusFail,
			Error:  fmt.Errorf("db connect failed: %w", err),
		}
	}
	defer conn.Close()

	// Verify schema was created.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return Result{
			Name:   "Database",
			Status: StatusFail,
			Error:  fmt.Errorf("database file not created"),
		}
	}

	return Result{
		Name:   "Database",
		Status: StatusPass,
	}
}
