package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Dotfiles struct {
	Name		string
	Filepath	string
}

type File struct {
	Name		string
	Src		string
	Dst		string
}

func New(homedir string, name string) (*Dotfiles, error) {
	dotfiles := &Dotfiles{
		Name:		name,
		Filepath:	path.Join(homedir, name),
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
			if ! StringInSlice(f.Name(), ignore) {
				result = append(result, f.Name())
			}
		} else {
			result = append(result, f.Name())
		}
	}

	return result, nil
}

func StringInSlice(a string, list []string) bool {
	for _, v := range list {
		if a == v {
			return true
		}
	}

	return false
}

func FileExist(f string) bool {
	_, err := os.Lstat(f)
	if err != nil {
		return false
	}

	return true
}

func isSymlink(f string) (bool, string, error) {
	t := false
	link := ""

	fi, err := os.Lstat(f)
	if err != nil {
		return t, link, err
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		t = true
		link, err = os.Readlink(f)
		if err != nil {
			return t, link, err
		}
	}

	return t, link, nil
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
