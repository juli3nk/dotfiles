# Dotfiles

```console
$ dotfiles -h
Usage of dotfiles:
  -dry-run=false: don't modify anything, just print commands
  -force=false: overwrite existing files
  -list=false: list currently managed dotfiles
  -name="Dotfiles": name of dotfiles repo
  -sync=false: update dotfile symlink
  -version=false: print version and exit
```

## Configuration

You can choose to create a configuration file inside your repository ``.dotfiles.yaml``.

To ignore files:

```
---
ignore:
  - '.git'
  - 'README.md'
```
