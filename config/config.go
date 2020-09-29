package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	}
	Cassandra struct {
		Host     string `yaml:"host"`
		KeySpace string `yaml:"keyspace"`
	}
	Mail struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Smtp     struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		}
	}
}

func GetConf() *Config {
	var config *Config
	yamlFile, err := ioutil.ReadFile("./config/app.yaml")
	if err != nil {
		fmt.Println("yamlFile get error", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}

	return config
}
