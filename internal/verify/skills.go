package verify

import (
	"context"
	"os"
	"path/filepath"
)

// VerifySkills checks the skills system.
func VerifySkills(ctx context.Context) Result {
	// Check that skills package exists.
	skillsDir := findSkillsDir()
	if skillsDir == "" {
		return Result{
			Name:   "Skills System",
			Status: StatusWarn,
			Details: "skills directory not found",
		}
	}

	// Check for key files.
	requiredFiles := []string{"discovery.go", "manager.go"}
	for _, file := range requiredFiles {
		path := filepath.Join(skillsDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return Result{
				Name:    "Skills System",
				Status:  StatusWarn,
				Details: file + " not found",
			}
		}
	}

	return Result{
		Name:   "Skills System",
		Status: StatusPass,
	}
}

func findSkillsDir() string {
	candidates := []string{
		"internal/skills",
		"../skills",
		"../../skills",
	}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return ""
}
