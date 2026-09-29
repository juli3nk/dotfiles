package profile

import (
	"fmt"
	"sort"

	"github.com/juli3nk/dotfiles/internal/helpers"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles from the config files",
		RunE:  runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	c, err := helpers.OpenComposed()
	if err != nil {
		return err
	}

	seen := make(map[string]bool)

	for _, cfg := range c.DotConfigs {
		for _, name := range cfg.ProfileNames() {
			seen[name] = true
		}
	}

	names := make([]string, 0, len(seen))

	for name := range seen {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}
