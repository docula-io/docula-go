package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd produces the root for the adr command tree.
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "adr",
		Short:   "ADR provides tooling for managing architecture decision records.",
		Long:    "ADR provides tooling for managing architecture decision records.",
		Version: "0.1.0",
	}

	rootCmd.AddCommand(initCmd())

	return rootCmd
}
