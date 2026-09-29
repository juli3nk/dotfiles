package repo

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

var repoAddOpts struct {
	force bool
}

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <git_url>",
		Short: "Add a repository and clone it",
		Args:  cobra.ExactArgs(2),
		RunE:  runAdd,
	}

	cmd.Flags().BoolVar(&repoAddOpts.force, "force", false, "overwrite existing repository directory")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	url := args[1]

	reg, err := helpers.OpenRegistry()
	if err != nil {
		return err
	}

	if reg.Exists(name) {
		return fmt.Errorf("repository '%s' is already registered", name)
	}

	target := reg.CloneDir(name)

	if exists, _ := dirHasFiles(target); exists {
		if !repoAddOpts.force {
			return fmt.Errorf("directory %s already exists, use --force to overwrite", target)
		}

		if err := os.RemoveAll(target); err != nil {
			return err
		}
	}

	if err := reg.EnsureCloneDir(); err != nil {
		return err
	}

	// #nosec G204
	// url and target are intentional user-provided git/clone arguments.
	clone := exec.Command("git", "clone", url, target)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr

	if err := clone.Run(); err != nil {
		return err
	}

	return reg.Add(name, url)
}
