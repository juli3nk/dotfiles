package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/juliengk/go-utils/filedir"
)

const VERSION = "v0.1.0"
const CONFIG_FILE = ".dotfiles.yaml"

var (
	flgForce   bool
	flgList    bool
	flgName    string
	flgSync    bool
	flgDryRun  bool
	flgVersion bool
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

	// Parse config file if exists
	_, err = os.Lstat(configfile)
	if err == nil {
		if err := config.Parse(configfile); err != nil {
			fmt.Println(err)
		}
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
			file := &File{
				Name: f,
				Src:  path.Join(dotfiles.Filepath, f),
				Dst:  path.Join(homedir, f),
			}

			if filedir.FileExists(file.Dst) == false {
				file.Symlink(flgDryRun)

				continue
			}

			sl, p, err := filedir.IsSymlink(file.Dst)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if sl == true && p == file.Src {
				continue
			}

			if (flgForce == false && sl == true && p != file.Src) || (flgForce == false && sl == false) {
				fmt.Printf("Skipping \"%s\", use -force to override\n", f)
			}

			if flgForce == true {
				file.Remove(flgDryRun)
				file.Symlink(flgDryRun)
			}
		}
	}
}
