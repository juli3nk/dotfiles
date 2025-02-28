package main

import (
	"context"
	"log"

	"dagger/dotfiles/internal/dagger"

	"github.com/juli3nk/go-utils/filedir"
)

// Lint Golang files
func (m *Dotfiles) LintGolang(
	ctx context.Context,
	checkOnlyModifiedFiles bool,
) (string, error) {
	opts := dagger.GoLintOpts{}

	if checkOnlyModifiedFiles {
		files := filedir.FilterFileByExtension(m.Git.ModifiedFiles, "go")
		if len(files) == 0 {
			log.Println("No modified Go files to lint.")
			return "", nil
		}
		opts.Filedir = files
	}

	return dag.Go(goVersion, m.Worktree).Lint(ctx, opts)
}

// Lint JSON files
func (m *Dotfiles) LintJsonFile(
	ctx context.Context,
	checkOnlyModifiedFiles bool,
) (string, error) {
	opts := dagger.JsonfileLintOpts{}

	if checkOnlyModifiedFiles {
		files := filedir.FilterFileByExtension(m.Git.ModifiedFiles, "json")
		if len(files) == 0 {
			log.Println("No modified JSON files to lint.")
			return "", nil
		}
		opts.Filedir = files
	}

	return dag.Jsonfile().Lint(ctx, m.Worktree, opts)
}
