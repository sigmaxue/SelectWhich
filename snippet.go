package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//解析yml文件
type Snippet struct {
	Name     string `yaml:"name"`
	Content  string `yaml:"content"`
	Comment  string `yaml:"comment"`
	CmdType  string `yaml:"cmdtype"`
}

type BaseInfo struct {
	Version     string `yaml:"version"`
	Snippets    []Snippet `yaml:"snippet"`
}

func (c *BaseInfo) GetConf(filename string) (*BaseInfo,  error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil 
}
