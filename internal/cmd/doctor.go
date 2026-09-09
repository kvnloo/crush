package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/crush/internal/verify"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run diagnostic checks on the Crush installation",
	Long: `The doctor command performs comprehensive health checks across the Crush codebase,
including configuration, agents, database, LSP, MCP, UI, hooks, skills, and testing infrastructure.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return runDoctor(ctx)
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(ctx context.Context) error {
	fmt.Println("🩺 Crush Doctor Report")
	fmt.Println("═════════════════════════════════════════════════")
	fmt.Println()

	checkers := verify.AllCheckers()
	results := make([]verify.Result, 0, len(checkers))

	// Run checks in a defined order.
	order := []string{
		"config",
		"agent",
		"db",
		"lsp",
		"mcp",
		"ui",
		"hooks",
		"skills",
		"testing",
	}

	for _, name := range order {
		checker, ok := checkers[name]
		if !ok {
			continue
		}
		result := checker(ctx)
		results = append(results, result)
		fmt.Println(result.String())
	}

	fmt.Println()

	// Summary.
	passed := 0
	failed := 0
	warned := 0
	for _, r := range results {
		switch r.Status {
		case verify.StatusPass:
			passed++
		case verify.StatusFail:
			failed++
		case verify.StatusWarn:
			warned++
		}
	}

	if failed > 0 {
		fmt.Printf("Overall: %d/%d checks passed, %d failed", passed, len(results), failed)
		if warned > 0 {
			fmt.Printf(", %d warnings", warned)
		}
		fmt.Println()
		os.Exit(1)
	}

	fmt.Printf("Overall: %d/%d checks passed", passed, len(results))
	if warned > 0 {
		fmt.Printf(" (%d warnings)", warned)
	}
	fmt.Println()

	return nil
}
