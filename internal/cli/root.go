package cli

import (
	"github.com/juli3nk/dotfiles/internal/cli/file"
	"github.com/juli3nk/dotfiles/internal/cli/profile"
	"github.com/juli3nk/dotfiles/internal/cli/repo"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dotfiles",
		Short: "Manage your dotfiles symlinks in your homedir",
	}

	cmd.AddCommand(file.NewCommand())
	cmd.AddCommand(profile.NewCommand())
	cmd.AddCommand(repo.NewCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}

func Execute() error {
	return NewCommand().Execute()
}
