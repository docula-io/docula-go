package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type initHandler func(ctx context.Context, path string) error

func initCmd(handler initHandler) *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init <path>",
		Short: "Sets up a directory as an ADR directory.",
		Long: "Sets up a directory as an ADR directory. " +
			"If the directory does not exist, then this command will create " +
			"the directory for the user.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			if err := handler(cmd.Context(), path); err != nil {
				return fmt.Errorf("init handler: %w", err)
			}

			return nil
		},
	}

	return initCmd
}
