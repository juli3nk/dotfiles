package main

import (
	"context"

	"dagger/dotfiles/internal/dagger"
)

// Lint commit messages
func (m *Dotfiles) LintCommitMsg(
	ctx context.Context,
	args []string,
) (string, error) {
	return dag.Commitlint().
		Lint(m.Worktree, dagger.CommitlintLintOpts{Args: args}).
		Stdout(ctx)
}
