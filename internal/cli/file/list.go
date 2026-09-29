package file

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/juli3nk/dotfiles/internal/dotfiles"
	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List currently managed dotfiles",
		RunE:  runList,
	}

	cmd.Flags().StringVar(&fileOpts.profile, "profile", "", "name of the profile to use")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	c, err := helpers.OpenComposed()
	if err != nil {
		return err
	}

	if err := validateProfile(c.Config, fileOpts.profile); err != nil {
		return err
	}

	homedir := os.Getenv("HOME")

	files, err := dotfiles.BuildComposedFiles(c.Dots, c.DotConfigs, c.Config, homedir, fileOpts.profile)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, "NAME\tSRC\tTYPE\tDST\tSTATE"); err != nil {
		return err
	}

	for _, f := range files {
		dst := strings.TrimPrefix(f.Dst, homedir)
		if len(dst) > 0 && dst[0] == '/' {
			dst = "~" + dst
		}

		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", f.Name, f.Repo, dstType(f), dst, f.State()); err != nil {
			return err
		}
	}

	return tw.Flush()
}
