package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

const VERSION = "v0.1.0"
const CONFIG_FILE = ".dotfiles.yaml"

var (
	flgForce	bool
	flgList		bool
	flgName		string
	flgSync		bool
	flgDryRun	bool
	flgVersion	bool
)

func init() {
	flag.BoolVar(&flgForce, "force", false, "overwrite existing files")
	flag.BoolVar(&flgList, "list", false, "list currently managed dotfiles")
	flag.StringVar(&flgName, "name", "Dotfiles", "name of dotfiles repo")
	flag.BoolVar(&flgSync, "sync", false, "update dotfile symlink")
	flag.BoolVar(&flgDryRun, "dry-run", false, "don't modify anything, just print commands")
	flag.BoolVar(&flgVersion, "version", false, "print version and exit")
	flag.Parse()
}

func main() {
	var config Config

	if flgVersion {
		fmt.Println(VERSION)
		return
	}

	name := flgName

	homedir := os.Getenv("HOME")

	dotfiles, err := New(homedir, name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get config file full path
	configfile := path.Join(dotfiles.Filepath, CONFIG_FILE)

	// Parse config file
	if err := config.Parse(configfile); err != nil {
		fmt.Println(err)
	}

	// Get the list of files
	files, err := GetFiles(dotfiles.Filepath, config.Ignore)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flgList {
		for _, f := range files {
			fmt.Println(f)
		}
	}

	if flgSync {
		for _, f := range files {
			overwrite := false

			src_f := path.Join(dotfiles.Filepath, f)
			dst_f := path.Join(homedir, f)

			fi, err := os.Lstat(dst_f)

			if err != nil {
				if flgDryRun {
					fmt.Printf("Creating symlink: %s\n", f)
				} else {
					os.Symlink(src_f, dst_f)
				}

				continue
			}

			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				dst_l, err := os.Readlink(dst_f)
				if err != nil {
					fmt.Println(err)
				}

				if dst_l != src_f {
					overwrite = true

					fmt.Printf("Skipping \"%s\", use -force to override\n", f)
				}
			} else {
				overwrite = true
			}

			if flgForce && overwrite {
				if flgDryRun {
					fmt.Printf("Removing %s\n", f)
					fmt.Printf("Creating symlink: %s\n", f)
				} else {
					os.Remove(dst_f)
					os.Symlink(src_f, dst_f)
				}
			}
		}
	}
}
