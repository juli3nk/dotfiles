package dotfiles

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/juli3nk/go-utils/filedir"
	"github.com/thoas/go-funk"

	"github.com/juli3nk/dotfiles/internal/config"
)

const ConfigFile = ".dotfiles.yml"

type Dotfiles struct {
	Name     string
	Filepath string
}

type File struct {
	Name string
	Src  string
	Dst  string
	Repo string
}

// State returns the real state of the destination:  linked, missing,
// mismatch (symlink to another target) or conflict (real file/directory).
func (f *File) State() string {
	fi, err := os.Lstat(f.Dst)
	if err != nil {
		if os.IsNotExist(err) {
			return "missing"
		}
		return "?"
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(f.Dst)
		if err != nil {
			return "?"
		}
		if target == f.Src {
			return "linked"
		}
		return "mismatch"
	}

	return "conflict"
}

func NewWithPath(name string, filepath string) (*Dotfiles, error) {
	d := &Dotfiles{
		Name:     name,
		Filepath: filepath,
	}

	if _, err := os.Stat(d.Filepath); os.IsNotExist(err) {
		return d, fmt.Errorf("dotfiles repository '%s' does not exist at '%s', run `dotfiles repo add`", name, filepath)
	}

	return d, nil
}

func (d *Dotfiles) ConfigFilePath() string {
	return path.Join(d.Filepath, ConfigFile)
}

func (d *Dotfiles) LoadConfig() (*config.Config, error) {
	return config.New(d.ConfigFilePath())
}

// MakeDirectories creates the directories declared in the merged
// configuration for the given profile.
func MakeDirectories(cfg *config.Config, homedir string, profile string, dryrun bool) {
	for _, dir := range cfg.GetDirectories(profile) {
		MakeDirectory(homedir, dir, dryrun)
	}
}

// BuildComposedFiles computes the target file set over the composed
// repositories. dotConfigs[i] is the parsed config of dots[i]. cfg is the
// merged configuration (used for the ignore set). Destinations are keyed and
// later repositories override earlier ones on conflicts.
func BuildComposedFiles(dots []*Dotfiles, dotConfigs []*config.Config, cfg *config.Config, homedir string, profile string) ([]*File, error) {
	filesByDst := make(map[string]*File)

	for i, d := range dots {
		links := dotConfigs[i].GetLinks(profile)
		ignore := cfg.GetIgnore(profile)

		var links2 []string

		for _, l := range links {
			s := strings.Split(l, ":")

			links2 = append(links2, s[0])

			f := &File{
				Name: s[0],
				Src:  path.Join(d.Filepath, s[0]),
				Repo: d.Name,
			}
			if len(s) == 2 {
				f.Dst = path.Join(homedir, s[1])
			} else {
				f.Dst = path.Join(homedir, s[0])
			}
			filesByDst[f.Dst] = f
		}

		found, err := GetFiles(d.Filepath, ignore, links2)
		if err != nil {
			return nil, err
		}

		for _, name := range found {
			f := &File{
				Name: name,
				Src:  path.Join(d.Filepath, name),
				Dst:  path.Join(homedir, name),
				Repo: d.Name,
			}
			filesByDst[f.Dst] = f
		}
	}

	files := make([]*File, 0, len(filesByDst))

	for _, f := range filesByDst {
		files = append(files, f)
	}

	return files, nil
}

func SyncFiles(files []*File, homedir string, dryrun bool, force bool, verbose bool) error {
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })

	if dryrun {
		fmt.Println("Dry run, no changes will be made")
	}

	var created, modified, skipped, ok int

	for _, file := range files {
		var verb, detail string

		switch file.State() {
		case "missing":
			if err := file.Symlink(dryrun); err != nil {
				return err
			}
			verb = "create"
			created++
		case "linked":
			verb = "ok"
			ok++
		case "mismatch":
			if force {
				if err := file.Remove(dryrun); err != nil {
					return err
				}
				if err := file.Symlink(dryrun); err != nil {
					return err
				}
				verb = "relink"
				modified++
			} else {
				verb = "skip"
				detail = "points elsewhere, use --force"
				skipped++
			}
		case "conflict":
			if force {
				if err := file.Remove(dryrun); err != nil {
					return err
				}
				if err := file.Symlink(dryrun); err != nil {
					return err
				}
				verb = "replace"
				modified++
			} else {
				verb = "skip"
				detail = "exists, use --force"
				skipped++
			}
		default:
			verb = "?"
		}

		if verb == "ok" && !verbose {
			continue
		}

		dst := strings.TrimPrefix(file.Dst, homedir)
		if len(dst) > 0 && dst[0] == '/' {
			dst = "~" + dst
		}

		fmt.Printf("%-8s %s", verb, dst)
		if detail != "" {
			fmt.Printf(" (%s)", detail)
		}
		fmt.Println()
	}

	fmt.Printf("%d created, %d modified, %d skipped, %d ok", created, modified, skipped, ok)
	if dryrun {
		fmt.Print(" (dry run)")
	}
	fmt.Println()

	return nil
}

func GetFiles(dotfilesdir string, ignore, links []string) ([]string, error) {
	var result []string

	files, err := os.ReadDir(dotfilesdir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.Name() == ConfigFile {
			continue
		}

		if f.Name() == SecretsDir || strings.HasSuffix(f.Name(), ageSuffix) {
			continue
		}

		if len(ignore) > 0 {
			if funk.Contains(ignore, f.Name()) {
				continue
			}
		}

		if len(links) > 0 {
			if funk.Contains(links, f.Name()) {
				continue
			}
		}

		result = append(result, f.Name())
	}

	return result, nil
}

func MakeDirectory(homedir string, dir config.Dir, dryrun bool) {
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

func (f *File) Remove(dryrun bool) error {
	if dryrun {
		return nil
	}

	return os.RemoveAll(f.Dst)
}

func (f *File) Symlink(dryrun bool) error {
	if dryrun {
		return nil
	}

	return os.Symlink(f.Src, f.Dst)
}
