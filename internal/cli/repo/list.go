package repo

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List registered repositories",
		RunE:    runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	reg, err := helpers.OpenRegistry()
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintln(tw, "NAME\tBRANCH\tAHEAD\tBEHIND\tSTATUS\tURL"); err != nil {
		return err
	}

	for _, e := range reg.Repos() {
		prefix := " "
		if reg.IsDefault(e.Name) {
			prefix = "*"
		}

		cloneDir := reg.CloneDir(e.Name)

		if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
			if _, err := fmt.Fprintf(tw, "%s %s\t-\t-\t-\tnot cloned\t%s\n", prefix, e.Name, e.URL); err != nil {
				return err
			}
			continue
		}

		s := gitInfo(cloneDir)

		ahead, behind := "-", "-"
		if s.hasUpstream {
			ahead = strconv.Itoa(s.ahead)
			behind = strconv.Itoa(s.behind)
		}

		status := "clean"
		if s.dirty > 0 {
			status = fmt.Sprintf("dirty(%d)", s.dirty)
		}

		if _, err := fmt.Fprintf(tw, "%s %s\t%s\t%s\t%s\t%s\t%s\n", prefix, e.Name, s.branch, ahead, behind, status, e.URL); err != nil {
			return err
		}
	}

	return tw.Flush()
}
