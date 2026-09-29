package dotfiles

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"filippo.io/age"
)

const (
	ageSuffix  = ".age"
	SecretsDir = "secrets"
)

type secretEntry struct {
	src  string
	dest string
}

// Identities loads the age identities from the given key file. Comments and
// blank lines are ignored, leading/trailing whitespace is trimmed. Returns an
// error when the file is unreadable or contains no valid identity.
func Identities(keyPath string) ([]age.Identity, error) {
	// #nosec G304
	// keyPath is the user-supplied age identity file expected to be read in full.
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read age identities from '%s': %w", keyPath, err)
	}

	var identities []age.Identity

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		id, err := age.ParseX25519Identity(line)
		if err != nil {
			return nil, fmt.Errorf("invalid age identity in '%s': %w", keyPath, err)
		}

		identities = append(identities, id)
	}

	if len(identities) == 0 {
		return nil, fmt.Errorf("no age identity found in '%s', generate one with `age-keygen -o %s`", keyPath, keyPath)
	}

	return identities, nil
}

// RefreshSecrets decrypts every *.age file of the composed repositories into
// <repo>/secrets/<relative path without .age> and removes stale plaintext
// files. With decrypt=false (dry-run) nothing is read or written and no key is
// required. When no .age file is found, RefreshSecrets is a no-op even without
// dry-run. A missing key or a decryption failure is a blocking error.
func RefreshSecrets(dots []*Dotfiles, keyPath string, decrypt bool) error {
	type repoSecrets struct {
		d       *Dotfiles
		secrets []secretEntry
	}

	var repos []repoSecrets
	found := false

	for _, d := range dots {
		secrets, err := collectSecrets(d.Filepath)
		if err != nil {
			return err
		}

		if len(secrets) > 0 {
			found = true
			repos = append(repos, repoSecrets{d: d, secrets: secrets})
		}
	}

	if !found || !decrypt {
		return nil
	}

	identities, err := Identities(keyPath)
	if err != nil {
		return err
	}

	for _, r := range repos {
		if err := decryptSecrets(r.d.Filepath, r.secrets, identities); err != nil {
			return err
		}
	}

	return nil
}

// collectSecrets walks the repository (recursively, skipping .git and the
// secrets directory itself) and returns the .age entries with their mirror
// destination path.
func collectSecrets(repoPath string) ([]secretEntry, error) {
	var secrets []secretEntry

	err := filepath.Walk(repoPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if p == repoPath {
				return nil
			}

			if base := path.Base(p); base == ".git" || base == SecretsDir {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(p, ageSuffix) {
			return nil
		}

		rel, err := filepath.Rel(repoPath, p)
		if err != nil {
			return err
		}

		secrets = append(secrets, secretEntry{
			src:  p,
			dest: path.Join(repoPath, SecretsDir, strings.TrimSuffix(rel, ageSuffix)),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// decryptSecrets writes the decrypted contents of each .age file to its mirror
// path with restrictive permissions (directories 0700, files 0600) and removes
// stale plaintext files that no longer have a .age source.
func decryptSecrets(repoPath string, secrets []secretEntry, identities []age.Identity) error {
	written := make(map[string]bool, len(secrets))

	for _, s := range secrets {
		plain, err := decryptFile(s.src, identities)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(path.Dir(s.dest), 0o700); err != nil {
			return err
		}
		if err := os.WriteFile(s.dest, plain, 0o600); err != nil {
			return err
		}
		written[s.dest] = true
	}

	return pruneSecrets(path.Join(repoPath, SecretsDir), written)
}

// decryptFile reads and decrypts a single .age file, returning its plaintext.
func decryptFile(src string, identities []age.Identity) (plain []byte, err error) {
	// #nosec G304
	// src is a .age file selected from the repository by collectSecrets.
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	r, err := age.Decrypt(f, identities...)
	if err != nil {
		return nil, fmt.Errorf("decrypting '%s': %w", src, err)
	}

	plain, err = io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decrypting '%s': %w", src, err)
	}

	return plain, nil
}

// pruneSecrets removes every file under the secrets directory that was not
// regenerated, then removes now-empty subdirectories.
func pruneSecrets(secretsDir string, written map[string]bool) error {
	var dirs []string

	err := filepath.Walk(secretsDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if p != secretsDir {
				dirs = append(dirs, p)
			}
			return nil
		}

		if !written[p] {
			// #nosec G122
			// The secrets directory is created and managed by this process; the
			// walked paths are deterministic mirrors of the .age source files.
			return os.Remove(p)
		}

		return nil
	})
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			if err := os.Remove(dir); err != nil {
				return err
			}
		}
	}

	return nil
}
