package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/juli3nk/dotfiles/version"
	"github.com/juli3nk/go-utils/filedir"
)

const CONFIG_FILE = ".dotfiles.yml"

var (
	flgDryRun  bool
	flgForce   bool
	flgList    bool
	flgName    string
	flgProfile string
	flgSync    bool
	flgVersion bool
)

func init() {
	flag.BoolVar(&flgDryRun, "dry-run", false, "don't modify anything, just print commands")
	flag.BoolVar(&flgForce, "force", false, "overwrite existing files")
	flag.BoolVar(&flgList, "list", false, "list currently managed dotfiles")
	flag.StringVar(&flgName, "name", "Dotfiles", "name of dotfiles repo")
	flag.StringVar(&flgProfile, "profile", "", "name of the profile to use")
	flag.BoolVar(&flgSync, "sync", false, "update dotfile symlink")
	flag.BoolVar(&flgVersion, "version", false, "print version and exit")
	flag.Parse()
}

func main() {
	if flgVersion {
		ver := version.New()
		ver.Show()
		return
	}

	homedir := os.Getenv("HOME")

	dotfiles, err := New(homedir, flgName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get config file full path
	configfile := path.Join(dotfiles.Filepath, CONFIG_FILE)

	config, err := NewConfig(configfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(flgProfile) > 0 {
		if err := config.existsProfile(flgProfile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Create directories
	directories := config.getDirectories(flgProfile)

	for _, dir := range directories {
		MakeDirectory(homedir, dir, flgDryRun)
	}

	// Get the list of files
	var files2 []*File
	var links2 []string

	ignore := config.getIgnore(flgProfile)
	links := config.getLinks(flgProfile)

	for _, l := range links {
		s := strings.Split(l, ":")

		links2 = append(links2, s[0])

		f := File{
			Name: s[0],
			Src:  path.Join(dotfiles.Filepath, s[0]),
		}
		if len(s) == 2 {
			f.Dst = path.Join(homedir, s[1])
		} else {
			f.Dst = path.Join(homedir, s[0])
		}
		files2 = append(files2, &f)
	}

	files, err := GetFiles(dotfiles.Filepath, ignore, links2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, file := range files {
		f := File{
			Name: file,
			Src:  path.Join(dotfiles.Filepath, file),
			Dst:  path.Join(homedir, file),
		}
		files2 = append(files2, &f)
	}

	// Execute
	if flgList {
		for _, file := range files2 {
			fmt.Println(file.Name)
		}

		os.Exit(0)
	}

	if flgSync {
		for _, file := range files2 {
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
				fmt.Printf("Skipping, use -force to override: %s\n", file.Name)
			}

			if flgForce == true {
				file.Remove(flgDryRun)
				file.Symlink(flgDryRun)
			}
		}
	}
}
