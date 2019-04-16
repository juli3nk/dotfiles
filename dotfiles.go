package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/filedir"
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

func GetFiles(dotfilesdir string, ignore, links []string) ([]string, error) {
	var result []string

	files, _ := ioutil.ReadDir(dotfilesdir)

	for _, f := range files {
		if f.Name() == CONFIG_FILE {
			continue
		}

		if len(ignore) > 0 {
			if utils.StringInSlice(f.Name(), ignore, false) {
				continue
			}
		}

		if len(links) > 0 {
			if utils.StringInSlice(f.Name(), links, false) {
				continue
			}
		}

		result = append(result, f.Name())
	}

	return result, nil
}

func MakeDirectory(homedir string, dir Dir, dryrun bool) {
	folder := path.Join(homedir, dir.Name)

	if dryrun {
		if !filedir.DirExists(folder) {
			fmt.Printf("Creating directory: %s\n", dir.Name)
		}
	} else {
		if err := filedir.CreateDirIfNotExist(folder, true, dir.Chmod); err != nil {
			fmt.Println(err)
		}
	}
}

func (f *File) Remove(dryrun bool) {
	mode := "file"

	fi, err := os.Lstat(f.Dst)
	if err != nil {
		fmt.Println(err)
	}

	if fi.IsDir() {
		mode = "directory"
	}

	if dryrun {
		fmt.Printf("Removing %s: %s\n", mode, f.Name)
		return
	}

	os.RemoveAll(f.Dst)
}

func (f *File) Symlink(dryrun bool) {
	if dryrun {
		fmt.Printf("Creating symlink: %s\n", f.Name)
	} else {
		os.Symlink(f.Src, f.Dst)
	}
}
