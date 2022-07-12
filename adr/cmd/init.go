package cmd

import (
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init <path>",
		Short: "Sets up a directory as an ADR directory.",
		Long: "Sets up a directory as an ADR directory. " +
			"If the directory does not exist, then this command will create " +
			"the directory for the user.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return initCmd
}
