package models

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ConfigData struct {
	Server struct {
		Name       string `yaml:"name"`
		Mode       string `yaml:"mode"`
		Production struct {
			Port         string `yaml:"port"`
			ReadTimeout  int    `yaml:"read-timeout"`
			WriteTimeout int    `yaml:"write-timeout"`
		} `yaml:"production"`
		Debug struct {
			Port         string `yaml:"port"`
			ReadTimeout  int    `yaml:"read-timeout"`
			WriteTimeout int    `yaml:"write-timeout"`
		} `yaml:"debug"`
	} `yaml:"server"`
}

type Config struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

func (c *ConfigData) LoadConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *ConfigData) GetConfig() *Config {
	config := &Config{}
	if c.Server.Mode == "debug" {
		config.Port = c.Server.Debug.Port
		config.ReadTimeout = c.Server.Debug.ReadTimeout
		config.WriteTimeout = c.Server.Debug.WriteTimeout
	} else {
		config.Port = c.Server.Production.Port
		config.ReadTimeout = c.Server.Production.ReadTimeout
		config.WriteTimeout = c.Server.Production.WriteTimeout
	}
	return config
}
