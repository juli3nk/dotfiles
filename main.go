package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const VERSION = "v0.1.0"

var (
	flgForce	bool
	flgList 	bool
	flgName 	string
	flgSync 	bool
	flgVersion	bool
)

func init() {
	flag.BoolVar(&flgForce, "force", false, "overwrite existing files")
	flag.BoolVar(&flgList, "list", false, "list currently managed dotfiles")
	flag.StringVar(&flgName, "name", "Dotfiles", "name of dotfiles repo")
	flag.BoolVar(&flgSync, "sync", false, "update dotfile symlink")
	flag.BoolVar(&flgVersion, "version", false, "print version and exit")
	flag.Parse()
}

func main() {
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

	if flgList {
		files, err := GetFiles(dotfiles.Filepath)

		if err == nil {
			for _, f := range files {
				fmt.Println(f)
			}
		}
	}

	if flgSync {
		var config Config
		configfile := path.Join(homedir, name, ".dotfiles.yaml")

		data, err := ioutil.ReadFile(configfile)

		if err != nil {
			fmt.Println(err)
		}

		if err := config.Parse(data); err != nil {
			fmt.Println(err)
		}

		files, err := GetFiles(dotfiles.Filepath)

		if err == nil {
			var files_b []string

			if len(config.Ignore) > 0 {
				for _, f := range files {
					if ! StringInSlice(f, config.Ignore) {
						files_b = append(files_b, f)
					}
				}
			} else {
				files_b = files
			}

			for _, f := range files_b {
				overwrite := false

				src_f := path.Join(dotfiles.Filepath, f)
				dst_f := path.Join(homedir, f)

				fi, err := os.Lstat(dst_f)

				if err != nil {
					os.Symlink(src_f, dst_f)
					return
				}

				if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
					dst_l, err := os.Readlink(dst_f)

					if err != nil {
						fmt.Println(err)
					}

					if dst_l != src_f {
						overwrite = true

						fmt.Println("Link %s does not point to %s", dst_f, src_f)
					}
				} else {
					overwrite = true
				}

				if flgForce && overwrite {
					os.Remove(dst_f)
					os.Symlink(src_f, dst_f)
				}
			}
		}
	}
}
