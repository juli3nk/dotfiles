package main

import (
	"fmt"
	"io/ioutil"
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

func NewConfig(filename string) (*Config, error) {
	if _, err := os.Lstat(filename); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) existsProfile(profileName string) error {
	if _, ok := c.Profiles[profileName]; ok {
		return nil
	}

	return fmt.Errorf("Profile '%s' does not exist", profileName)
}

func (c *Config) getDirectories(profileName string) []Dir {
	var result []Dir

	for _, dir := range c.Common.Dirs {
		if dir.Chmod == 0000 {
			dir.Chmod = 0775
		}

		result = append(result, dir)
	}

	if len(profileName) == 0 {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			for _, dir := range c.Templates[tpl].Dirs {
				if dir.Chmod == 0000 {
					dir.Chmod = 0775
				}

				result = append(result, dir)
			}
		}
	}

	if len(c.Profiles[profileName].Dirs) == 0 {
		return result
	}
	for _, dir := range c.Profiles[profileName].Dirs {
		if dir.Chmod == 0000 {
			dir.Chmod = 0775
		}

		result = append(result, dir)
	}

	return result
}

func (c *Config) getIgnore(profileName string) []string {
	var result []string

	for _, ignore := range c.Common.Ignore {
		result = append(result, ignore)
	}

	if profileName == "" {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			for _, ignore := range c.Templates[tpl].Ignore {
				result = append(result, ignore)
			}
		}
	}

	if len(c.Profiles[profileName].Ignore) == 0 {
		return result
	}
	for _, ignore := range c.Profiles[profileName].Ignore {
		result = append(result, ignore)
	}

	return result
}

func (c *Config) getLinks(profileName string) []string {
	var result []string

	for _, link := range c.Common.Links {
		result = append(result, link)
	}

	if profileName == "" {
		return result
	}

	if len(c.Profiles[profileName].Include) > 0 {
		for _, tpl := range c.Profiles[profileName].Include {
			for _, link := range c.Templates[tpl].Links {
				result = append(result, link)
			}
		}
	}

	if len(c.Profiles[profileName].Links) == 0 {
		return result
	}
	for _, link := range c.Profiles[profileName].Links {
		result = append(result, link)
	}

	return result
}
