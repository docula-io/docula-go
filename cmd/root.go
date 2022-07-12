package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	adrCmd "github.com/docula-io/docula/adr/cmd"
)

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "docula [command]",
		Short:   "Docula provides tooling for various documentation types.",
		Long:    "Docula provides tooling for various documentation types.",
		Version: "0.1.0",
	}

	rootCmd.AddCommand(adrCmd.RootCmd())

	return rootCmd
}

// Execute acts as the main entry for the docules command cli. It loads
// the root command which in turn loads the sub commands for use.
func Execute(ctx context.Context) error {
	if err := rootCmd().ExecuteContext(ctx); err != nil {
		return fmt.Errorf("executing root command: %w", err)
	}

	return nil
}
