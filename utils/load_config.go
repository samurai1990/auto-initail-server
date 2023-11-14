package utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type YamlInfo struct {
	IP       string `yaml:"ip"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	NewPort  int    `yaml:"new_port"`
	Password string `yaml:"password"`
}

type Config struct {
	Path  string
	Yamls *[]YamlInfo
}

func Newconfig(path string) *Config {
	return &Config{
		Path: path,
	}
}

func (c *Config) GetConf() error {

	yamlFile, err := os.ReadFile(c.Path)
	if err != nil {
		return fmt.Errorf("yamlFile.Get err   #%v ", err)
	}

	var configs []YamlInfo
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	c.Yamls = &configs

	return nil
}
