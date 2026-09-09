package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/crush/internal/verify"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [feature]",
	Short: "Run deep verification on specific features",
	Long: `The verify command runs deep verification checks on specific features.
Available features: config, agent, db, lsp, mcp, ui, hooks, skills, testing, all`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		feature := "all"
		if len(args) > 0 {
			feature = args[0]
		}
		return runVerify(ctx, feature)
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func runVerify(ctx context.Context, feature string) error {
	checkers := verify.AllCheckers()

	if feature == "all" {
		fmt.Println("Running all feature verifications...")
		fmt.Println()

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
			fmt.Println(result.String())
		}

		return nil
	}

	checker, ok := checkers[feature]
	if !ok {
		return fmt.Errorf("unknown feature: %s\nAvailable: config, agent, db, lsp, mcp, ui, hooks, skills, testing, all", feature)
	}

	result := checker(ctx)
	fmt.Println(result.String())

	if result.Status == verify.StatusFail {
		return fmt.Errorf("verification failed")
	}

	return nil
}
