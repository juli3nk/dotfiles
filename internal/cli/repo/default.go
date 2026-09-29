package repo

import (
	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

func newDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default <name>",
		Short: "Set the default repository (the base/socle)",
		Args:  cobra.ExactArgs(1),
		RunE:  runDefault,
	}

	return cmd
}

func runDefault(cmd *cobra.Command, args []string) error {
	reg, err := helpers.OpenRegistry()
	if err != nil {
		return err
	}

	return reg.SetDefault(args[0])
}
