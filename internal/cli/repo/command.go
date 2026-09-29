package repo

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage registered dotfiles repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		newListCommand(),
		newAddCommand(),
		newRemoveCommand(),
		newDefaultCommand(),
	)

	return cmd
}

func dirHasFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return len(entries) > 0, nil
}

type gitStatus struct {
	branch      string
	ahead       int
	behind      int
	dirty       int
	hasUpstream bool
}

func gitRun(dir string, args ...string) (string, error) {
	// #nosec G204
	// args are hardcoded git subcommands passed by the callers of this helper.
	c := exec.Command("git", args...)
	c.Dir = dir

	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func gitInfo(dir string) *gitStatus {
	s := &gitStatus{branch: "?"}

	// Fetch the latest refs, but never fail the whole listing if it fails
	// (offline, no remote): fall back to the local refs.
	_, _ = gitRun(dir, "fetch")

	out, err := gitRun(dir, "status", "--porcelain=v2", "--branch")
	if err != nil {
		return s
	}

	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "# branch.head "):
			s.branch = strings.TrimPrefix(line, "# branch.head ")
		case strings.HasPrefix(line, "# branch.upstream "):
			s.hasUpstream = true
		case strings.HasPrefix(line, "# branch.ab "):
			fields := strings.Fields(strings.TrimPrefix(line, "# branch.ab "))
			if len(fields) == 2 {
				s.ahead, _ = strconv.Atoi(fields[0])
				s.behind, _ = strconv.Atoi(fields[1])
				if s.behind < 0 {
					s.behind = -s.behind
				}
			}
		case strings.HasPrefix(line, "#"):
			continue
		default:
			s.dirty++
		}
	}

	if s.branch == "(detached)" {
		if hash, err := gitRun(dir, "rev-parse", "--short", "HEAD"); err == nil {
			s.branch = hash
		}
	}

	return s
}
