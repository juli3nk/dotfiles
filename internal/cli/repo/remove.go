package repo

import (
	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <name>",
		Aliases: []string{"rm"},
		Short:   "Unregister a repository (cloned directory is left intact)",
		Args:    cobra.ExactArgs(1),
		RunE:    runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	reg, err := helpers.OpenRegistry()
	if err != nil {
		return err
	}

	return reg.Remove(args[0])
}
