# Dotfiles

`dotfiles` is a tool to manage your dot-files symlinks in homedir.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
- [Usage](#usage)
- [Repositories & XDG](#repositories--xdg)
- [Example Configuration File](#example-configuration-file)
- [Secrets](#secrets)
- [Configuration Options](#configuration-options)

<!-- tocstop -->

## Installation

### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/juli3nk/dotfiles/releases).

### Via Go

```console
$ go get github.com/juli3nk/dotfiles
```

## Usage

Dotfiles are managed through multiple git repositories registered in a local
registry. One repository is the *default* (the base/socle); the others are
layers merged on top of it, later repositories overriding earlier ones on
conflicts.

```console
$ dotfiles repo add work https://github.com/juli3nk/dot-files-work.git   # first add becomes the default
$ dotfiles repo list
$ dotfiles profile ls --profile work

$ dotfiles file list            # list managed dotfiles
$ dotfiles file sync --dry-run  # preview symlinks and dirs
$ dotfiles file sync            # create links
$ dotfiles file sync --force    # overwrite existing files

$ dotfiles version
```

## Repositories & XDG

The registry and the cloned repositories follow the
[XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/):

| What                 | Path (default)                              |
|----------------------|---------------------------------------------|
| Registry file        | `$XDG_CONFIG_HOME/dotfiles/repos.yml` (`~/.config/dotfiles/repos.yml`) |
| Repository clones    | `$XDG_DATA_HOME/dotfiles/repos/` (`~/.local/share/dotfiles/repos/`) |

If `XDG_CONFIG_HOME`/`XDG_DATA_HOME` are unset (or empty), the `~/.config` and
`~/.local/share` fallbacks are used.

## Example Configuration File

Create a configuration file `.dotfiles.yml` inside your dot-files repository.

```yaml
---
common:
  dirs:
    - name: '.config'
    - name: '.shell_custom.d'
  ignore:
    - '.git'
    - '.gitignore'
    - 'README.md'
    - '.config'
    - '.shell_custom.d'

templates:
  template1:
    ignore:
      - '.i3'
      - '.Xresources'

profiles:
  me:
    links:
      - '.config/terminator'
      - '.shell_custom.d/dockerfunc.sh'
  nox:
    include:
      - 'template1'
```

## Secrets

Encrypted files can be committed to your repositories with the
[age encryption format](https://age-encryption.org). A file ending with
`.age` is never deployed as-is: `dotfiles file sync` decrypts it into the
`secrets/` directory of the repository clone, mirroring its relative path.

```console
# one-time setup: generate an identity file
$ age-keygen -o ~/.config/age/keys.txt

# encrypt a file inside your dot-files repository and commit the .age file
$ cd <repository-clone>
$ age -e -r <AGE-PUBLIC-KEY> -o .ssh/config.age .ssh/config
$ git add .ssh/config.age
```

Once decrypted, link the plaintext into your home dir with the `secrets/`
prefix:

```yaml
common:
  ignore:
    - .ssh
profiles:
  work:
    links:
      - 'secrets/.ssh/config:.ssh/config'
      - 'secrets/.ssh/id_ed25519:.ssh/id_ed25519'
```

The `secrets/` directory holds plaintext only, so add it to the `.gitignore`
of every repository that declares encrypted files:

```
secrets/
```

Notes:

- `secrets/` and `*.age` files are ignored by the automatic top-level
  discovery; deployment happens exclusively through explicit `links`.
- Decrypted files are written with `0600` permissions (`0700` for
  directories); stale plaintext files are pruned at each sync.
- The identity file is looked up at `$XDG_CONFIG_HOME/age/keys.txt`
  (`~/.config/age/keys.txt`), can be overridden with the `--age-key` flag of
  `dotfiles file sync` or the `DOTFILES_AGE_KEY` environment variable, and may
  contain several identities (one per line, comments with `#`).
- A missing key, an unknown identity or a corrupted `.age` file aborts the
  sync before any symlink is touched.

## Configuration Options
### Common

Configurations common to all profiles.

### Templates

Allows to create redundant configurations to be included in profiles.

### Profiles

### Dirs

Create directory in your home dir.

### Links

Create symlink for specific file.

### Ignore

Ignore the soft link creation.

### Include

Include template in a profile.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
