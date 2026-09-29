package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Dir struct {
	Name  string      `yaml:"name"`
	User  string      `yaml:"user,omitempty"`
	Group string      `yaml:"group,omitempty"`
	Chmod os.FileMode `yaml:"chmod,omitempty"`
}

type Option struct {
	Include []string `yaml:"include,omitempty"`
	Dirs    []Dir    `yaml:"dirs,omitempty"`
	Ignore  []string `yaml:"ignore,omitempty"`
	Links   []string `yaml:"links,omitempty"`
}

type Config struct {
	Common    Option            `yaml:"common,omitempty"`
	Templates map[string]Option `yaml:"templates,omitempty"`
	Profiles  map[string]Option `yaml:"profiles,omitempty"`
}

func New(filename string) (*Config, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, err
	}

	// #nosec G304
	// filename is the configuration file the caller explicitly asks to load.
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := new(Config)

	if err = yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}

// Merge appends the option sets of other into c. Common dirs, ignore and
// links are concatenated. Templates and profiles maps are merged per key: an
// existing key's fields are appended, a new key is added as-is.
func (c *Config) Merge(other *Config) {
	c.Common.Dirs = append(c.Common.Dirs, other.Common.Dirs...)
	c.Common.Ignore = append(c.Common.Ignore, other.Common.Ignore...)
	c.Common.Links = append(c.Common.Links, other.Common.Links...)

	c.Templates = mergeOptionsMap(c.Templates, other.Templates)
	c.Profiles = mergeOptionsMap(c.Profiles, other.Profiles)
}

func mergeOptionsMap(dst, src map[string]Option) map[string]Option {
	if len(src) == 0 {
		return dst
	}

	if dst == nil {
		dst = make(map[string]Option, len(src))
	}

	for name, opt := range src {
		if existing, ok := dst[name]; ok {
			existing.Include = append(existing.Include, opt.Include...)
			existing.Dirs = append(existing.Dirs, opt.Dirs...)
			existing.Ignore = append(existing.Ignore, opt.Ignore...)
			existing.Links = append(existing.Links, opt.Links...)
			dst[name] = existing
			continue
		}

		dst[name] = opt
	}

	return dst
}

func (c *Config) ExistsProfile(profileName string) error {
	if _, ok := c.Profiles[profileName]; ok {
		return nil
	}

	return fmt.Errorf("profile '%s' does not exist", profileName)
}

func (c *Config) ProfileNames() []string {
	names := make([]string, 0, len(c.Profiles))

	for name := range c.Profiles {
		names = append(names, name)
	}

	return names
}

func (c *Config) GetDirectories(profileName string) []Dir {
	var result []Dir

	for _, dir := range c.Common.Dirs {
		if dir.Chmod == 0o000 {
			dir.Chmod = 0o775
		}
		result = append(result, dir)
	}

	if len(profileName) == 0 {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			for _, dir := range c.Templates[tpl].Dirs {
				if dir.Chmod == 0o000 {
					dir.Chmod = 0o775
				}
				result = append(result, dir)
			}
		}
	}

	if len(c.Profiles[profileName].Dirs) == 0 {
		return result
	}

	for _, dir := range c.Profiles[profileName].Dirs {
		if dir.Chmod == 0o000 {
			dir.Chmod = 0o775
		}
		result = append(result, dir)
	}

	return result
}

func (c *Config) GetIgnore(profileName string) []string {
	var result []string

	result = append(result, c.Common.Ignore...)

	if profileName == "" {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			result = append(result, c.Templates[tpl].Ignore...)
		}
	}

	if len(c.Profiles[profileName].Ignore) == 0 {
		return result
	}

	result = append(result, c.Profiles[profileName].Ignore...)

	return result
}

func (c *Config) GetLinks(profileName string) []string {
	var result []string

	result = append(result, c.Common.Links...)

	if profileName == "" {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			result = append(result, c.Templates[tpl].Links...)
		}
	}

	if len(c.Profiles[profileName].Links) == 0 {
		return result
	}

	result = append(result, c.Profiles[profileName].Links...)

	return result
}
