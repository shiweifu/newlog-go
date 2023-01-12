package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BlogPath  string `yaml:"blog_path"`
	BlogTitle string `yaml:"blog_title"`
	Port      int    `yaml:"port"`
}

func (c *Config) GetPort() int {
	if c.Port == 0 {
		return 8080
	}
	return c.Port
}

func readConf(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
