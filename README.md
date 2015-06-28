# Dotfiles

```console
$ dotfiles -h
Usage of dotfiles:
  -force=false: overwrite existing files
  -list=false: list currently managed dotfiles
  -name="Dotfiles": name of your dotfiles repo
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
