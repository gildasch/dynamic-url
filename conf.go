package main

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type conf struct {
	Movies map[string]movieConf
}

type movieConf struct {
	Movie     string
	Subtitles string
	Script    string
}

func parseConf(s string) (conf, error) {
	var c conf
	err := yaml.Unmarshal([]byte(s), &c)
	return c, err
}

func parseConfFile(path string) (conf, error) {
	f, err := os.Open(path)
	if err != nil {
		return conf{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return conf{}, err
	}

	return parseConf(string(b))
}
