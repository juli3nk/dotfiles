package file

import (
	"os"

	"github.com/juli3nk/dotfiles/internal/dotfiles"
	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

type fileSyncOptions struct {
	profile string
	dryRun  bool
	force   bool
	verbose bool
	ageKey  string
}

var fileSyncOpts fileSyncOptions

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Update dotfile symlinks",
		RunE:  runSync,
	}

	cmd.Flags().StringVar(&fileSyncOpts.profile, "profile", "", "name of the profile to use")
	cmd.Flags().BoolVar(&fileSyncOpts.dryRun, "dry-run", false, "don't modify anything, just print commands")
	cmd.Flags().BoolVarP(&fileSyncOpts.force, "force", "f", false, "overwrite existing files")
	cmd.Flags().BoolVarP(&fileSyncOpts.verbose, "verbose", "v", false, "show already linked files")
	cmd.Flags().StringVar(&fileSyncOpts.ageKey, "age-key", "", "path to the age identity file used to decrypt .age files")

	return cmd
}

func runSync(cmd *cobra.Command, args []string) error {
	homedir := os.Getenv("HOME")

	c, err := helpers.OpenComposed()
	if err != nil {
		return err
	}

	if err := validateProfile(c.Config, fileSyncOpts.profile); err != nil {
		return err
	}

	if err := dotfiles.RefreshSecrets(c.Dots, ageKeyPath(homedir, fileSyncOpts.ageKey), !fileSyncOpts.dryRun); err != nil {
		return err
	}

	dotfiles.MakeDirectories(c.Config, homedir, fileSyncOpts.profile, fileSyncOpts.dryRun)

	files, err := dotfiles.BuildComposedFiles(c.Dots, c.DotConfigs, c.Config, homedir, fileSyncOpts.profile)
	if err != nil {
		return err
	}

	return dotfiles.SyncFiles(files, homedir, fileSyncOpts.dryRun, fileSyncOpts.force, fileSyncOpts.verbose)
}
