# Dotfiles

`dotfiles` is a tool to manage your dot-files symlinks in homedir.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
- [Usage](#usage)
- [Example Configuration File](#example-configuration-file)
- [Configuration Options](#configuration-options)

<!-- tocstop -->

## Installation

### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/juliengk/dotfiles/releases).

### Via Go

```console
$ go get github.com/juli3nk/dotfiles
```

## Usage

```console
$ dotfiles -h
Usage of bin/dotfiles:
  -dry-run
        don't modify anything, just print commands
  -force
        overwrite existing files
  -list
        list currently managed dotfiles
  -name string
        name of dotfiles repo (default "Dotfiles")
  -profile string
        name of the profile to use
  -sync
        update dotfile symlink
  -version
        print version and exit
```

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
