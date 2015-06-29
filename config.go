package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Ignore []string
}

func (c *Config) Parse(configfile string) error {
	data, err := ioutil.ReadFile(configfile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)

	return err
}
