package registry

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const (
	configFile = "repos.yml"
)

// configBaseDir returns the base XDG config directory, falling back to
// ~/.config when XDG_CONFIG_HOME is unset or empty.
func configBaseDir(homedir string) string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}

	return path.Join(homedir, ".config")
}

// dataBaseDir returns the base XDG data directory, falling back to
// ~/.local/share when XDG_DATA_HOME is unset or empty.
func dataBaseDir(homedir string) string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return dir
	}

	return path.Join(homedir, ".local", "share")
}

type Entry struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Default string  `yaml:"default,omitempty"`
	Repos   []Entry `yaml:"repos,omitempty"`
}

type Registry struct {
	configPath string
	reposPath  string
	config     *Config
}

func New(homedir string) (*Registry, error) {
	r := &Registry{
		configPath: path.Join(configBaseDir(homedir), "dotfiles", configFile),
		reposPath:  path.Join(dataBaseDir(homedir), "dotfiles", "repos"),
	}

	cfg, err := loadConfig(r.configPath)
	if err != nil {
		return nil, err
	}

	r.config = cfg

	return r, nil
}

func loadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Lstat(configPath); err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	// #nosec G304
	// configPath is the registry file managed by this package.
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (r *Registry) save() error {
	data, err := yaml.Marshal(r.config)
	if err != nil {
		return err
	}

	dir := path.Dir(r.configPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return err
	}
	return os.WriteFile(r.configPath, data, 0o600)
}

func (r *Registry) CloneDir(name string) string {
	return path.Join(r.reposPath, name)
}

func (r *Registry) Default() string {
	return r.config.Default
}

func (r *Registry) IsDefault(name string) bool {
	return name != "" && name == r.config.Default
}

func (r *Registry) Repos() []Entry {
	return r.config.Repos
}

func (r *Registry) Exists(name string) bool {
	for _, e := range r.config.Repos {
		if e.Name == name {
			return true
		}
	}

	return false
}

func (r *Registry) Find(name string) (Entry, error) {
	for _, e := range r.config.Repos {
		if e.Name == name {
			return e, nil
		}
	}

	return Entry{}, fmt.Errorf("repository '%s' is not registered", name)
}

func (r *Registry) Add(name, url string) error {
	if r.Exists(name) {
		return fmt.Errorf("repository '%s' is already registered", name)
	}

	r.config.Repos = append(r.config.Repos, Entry{Name: name, URL: url})

	if r.config.Default == "" {
		r.config.Default = name
	}

	return r.save()
}

func (r *Registry) Remove(name string) error {
	var repos []Entry
	found := false

	for _, e := range r.config.Repos {
		if e.Name == name {
			found = true
			continue
		}

		repos = append(repos, e)
	}

	if !found {
		return fmt.Errorf("repository '%s' is not registered", name)
	}

	r.config.Repos = repos

	if r.config.Default == name {
		r.config.Default = ""
	}

	return r.save()
}

func (r *Registry) SetDefault(name string) error {
	if !r.Exists(name) {
		return fmt.Errorf("repository '%s' is not registered", name)
	}

	r.config.Default = name

	return r.save()
}

// Composed returns the repositories in composition order: the default
// repository (the base/socle) first, then the remaining repositories in
// registration order. Later repositories override the earlier ones.
func (r *Registry) Composed() ([]Entry, error) {
	if len(r.config.Repos) == 0 {
		return nil, nil
	}

	if r.config.Default == "" {
		return nil, fmt.Errorf("no default repository set, run `dotfiles repo default <name>`")
	}

	composed := make([]Entry, 0, len(r.config.Repos))

	var defaultRepo *Entry

	for i := range r.config.Repos {
		e := r.config.Repos[i]
		if e.Name == r.config.Default {
			defaultRepo = &r.config.Repos[i]
			continue
		}

		composed = append(composed, e)
	}

	if defaultRepo == nil {
		return nil, fmt.Errorf("default repository '%s' is not registered", r.config.Default)
	}

	composed = append([]Entry{*defaultRepo}, composed...)

	return composed, nil
}

// EnsureCloneDir creates the base directory holding cloned repositories.
func (r *Registry) EnsureCloneDir() error {
	return os.MkdirAll(r.reposPath, 0o750)
}
