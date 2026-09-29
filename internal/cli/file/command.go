package file

import (
	"os"
	"path"

	"github.com/juli3nk/dotfiles/internal/config"
	"github.com/juli3nk/dotfiles/internal/dotfiles"
	"github.com/spf13/cobra"
)

type fileOptions struct {
	profile string
}

var fileOpts fileOptions

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Manage dotfiles files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		newListCommand(),
		newSyncCommand(),
	)

	return cmd
}

// dstType returns the real filesystem type at the destination, like `ls -l`.
func dstType(f *dotfiles.File) string {
	// #nosec G703
	// f.Dst comes from the configured dotfiles destination mapping; listing the
	// real filesystem state at that path is the intended behaviour of the command.
	fi, err := os.Lstat(f.Dst)
	if err != nil {
		if os.IsNotExist(err) {
			return "-"
		}
		return "?"
	}

	switch {
	case fi.Mode()&os.ModeSymlink != 0:
		return "link"
	case fi.IsDir():
		return "dir"
	default:
		return "file"
	}
}

// ageKeyPath resolves the age identity file: --age-key flag, DOTFILES_AGE_KEY
// environment variable, then $XDG_CONFIG_HOME/age/keys.txt (or
// ~/.config/age/keys.txt).
func ageKeyPath(homedir string, flag string) string {
	if flag != "" {
		return flag
	}

	if env := os.Getenv("DOTFILES_AGE_KEY"); env != "" {
		return env
	}

	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = path.Join(homedir, ".config")
	}

	return path.Join(dir, "age", "keys.txt")
}

func validateProfile(cfg *config.Config, profile string) error {
	if len(profile) == 0 {
		return nil
	}

	return cfg.ExistsProfile(profile)
}
