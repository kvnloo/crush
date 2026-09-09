package verify

import (
	"context"
	"fmt"
)

// Status represents the result of a verification check.
type Status string

const (
	// StatusPass indicates the check passed.
	StatusPass Status = "PASS"
	// StatusFail indicates the check failed.
	StatusFail Status = "FAIL"
	// StatusWarn indicates the check passed with warnings.
	StatusWarn Status = "WARN"
	// StatusSkip indicates the check was skipped.
	StatusSkip Status = "SKIP"
)

// Result represents the outcome of a verification check.
type Result struct {
	Name    string
	Status  Status
	Error   error
	Details string
}

// String formats the result as a human-readable string.
func (r Result) String() string {
	var statusIcon string
	switch r.Status {
	case StatusPass:
		statusIcon = "✓"
	case StatusFail:
		statusIcon = "✗"
	case StatusWarn:
		statusIcon = "⚠"
	case StatusSkip:
		statusIcon = "○"
	}

	msg := fmt.Sprintf("%s %-18s %s", statusIcon, r.Name, r.Status)
	if r.Error != nil {
		msg += fmt.Sprintf(": %v", r.Error)
	}
	if r.Details != "" {
		msg += fmt.Sprintf(" (%s)", r.Details)
	}
	return msg
}

// Checker is a function that performs a verification check.
type Checker func(context.Context) Result

// AllCheckers returns all available verification checkers.
func AllCheckers() map[string]Checker {
	return map[string]Checker{
		"config":  VerifyConfig,
		"agent":   VerifyAgent,
		"db":      VerifyDatabase,
		"lsp":     VerifyLSP,
		"mcp":     VerifyMCP,
		"ui":      VerifyUI,
		"hooks":   VerifyHooks,
		"skills":  VerifySkills,
		"testing": VerifyTesting,
	}
}
