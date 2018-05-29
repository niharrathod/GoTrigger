package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the config.yaml
type Config struct {
	Server struct {
		Port        int    `yaml:"port"`
		ContextPath string `yaml:"contextPath"`
		LogPath     string `yaml:"logPath"`
	} `yaml:"server"`
	Scheduler struct {
		MinScheduleWindowSeconds int `yaml:"minScheduleWindowSeconds"`
		MaxScheduleWindowHours   int `yaml:"maxScheduleWindowHours"`
	} `yaml:"scheduler"`
}

var config Config

// LoadConfig loads the config from given yaml path.
func LoadConfig(path string) {
	ymlpath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	ymlConfig, err := ioutil.ReadFile(ymlpath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(ymlConfig, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
