package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/juliengk/go-utils"
)

type Dotfiles struct {
	Name     string
	Filepath string
}

type File struct {
	Name string
	Src  string
	Dst  string
}

func New(homedir string, name string) (*Dotfiles, error) {
	dotfiles := &Dotfiles{
		Name:     name,
		Filepath: path.Join(homedir, name),
	}

	if _, err := os.Stat(dotfiles.Filepath); os.IsNotExist(err) {
		return dotfiles, fmt.Errorf("Dotfiles '%s' directory does not exist, please specify a dotfiles directory with -name.", name)
	}

	return dotfiles, nil
}

func GetFiles(dotfilesdir string, ignore []string) ([]string, error) {
	var result []string

	files, _ := ioutil.ReadDir(dotfilesdir)

	// Remove the files to ignore from the list
	for _, f := range files {
		if f.Name() == ".dotfiles.yaml" {
			continue
		} else if len(ignore) > 0 {
			if !utils.StringInSlice(f.Name(), ignore, false) {
				result = append(result, f.Name())
			}
		} else {
			result = append(result, f.Name())
		}
	}

	return result, nil
}

func (f *File) Symlink(dryrun bool) {
	if dryrun {
		fmt.Printf("Creating symlink: %s\n", f.Name)
	} else {
		os.Symlink(f.Src, f.Dst)
	}
}

func (f *File) Remove(dryrun bool) {
	if dryrun {
		fmt.Printf("Removing %s\n", f.Name)
	} else {
		os.Remove(f.Dst)
	}
}
