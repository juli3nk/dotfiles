package main

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	Ignore []string
}

func (c *Config) Parse(data []byte) error {
	err := yaml.Unmarshal(data, c)

	return err
}
