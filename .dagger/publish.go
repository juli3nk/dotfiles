package main

import (
	"context"
	"fmt"
	"net/url"

	"dagger/dotfiles/internal/dagger"
)

// Publish the binaries to Github Release
func (m *Dotfiles) Publish(
	ctx context.Context,
	Token *dagger.Secret,
) error {
	if len(m.Files) == 0 {
		return fmt.Errorf("error: build files first")
	}

	url, err := url.Parse(appSourceUrl)
	if err != nil {
		return fmt.Errorf("error: failed to parse app source url: %w", err)
	}

	opts := dagger.GhReleaseCreateOpts{
		GenerateNotes: true,
		Files:         m.Files,
		VerifyTag:     true,
		Repo:          url.Path,
		Token:         Token,
	}

	err = dag.Gh().Release().Create(ctx, m.Git.Tag, m.Git.Tag, opts)

	return err
}
