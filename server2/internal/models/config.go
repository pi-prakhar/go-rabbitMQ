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
			Port           string `yaml:"port"`
			ReadTimeout    int    `yaml:"read-timeout"`
			WriteTimeout   int    `yaml:"write-timeout"`
			HandlerTimeout int    `yaml:"hander-timeout"`
		} `yaml:"production"`
		Debug struct {
			Port           string `yaml:"port"`
			ReadTimeout    int    `yaml:"read-timeout"`
			WriteTimeout   int    `yaml:"write-timeout"`
			HandlerTimeout int    `yaml:"hander-timeout"`
		} `yaml:"debug"`
	} `yaml:"server"`
	RabbitMQ struct {
		Host       string `yaml:"host"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		RetryCount int    `yaml:"retry-count"`
		RetrySleep int    `yaml:"retry-sleep"`
	} `yaml:"rabbitmq"`
}

type Config struct {
	Port           string
	ReadTimeout    int
	WriteTimeout   int
	HandlerTimeout int
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

	modeFromEnv := os.Getenv("MODE")
	if modeFromEnv != "" {
		c.Server.Mode = modeFromEnv
	}

	rabbitMQHostFromEnv := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHostFromEnv != "" {
		c.RabbitMQ.Host = rabbitMQHostFromEnv
	}

	rabbitMQUserFromEnv := os.Getenv("RABBITMQ_USER")
	if rabbitMQUserFromEnv != "" {
		c.RabbitMQ.User = rabbitMQUserFromEnv
	}

	rabbitMQPasswordFromEnv := os.Getenv("RABBITMQ_PASS")
	if rabbitMQUserFromEnv != "" {
		c.RabbitMQ.Password = rabbitMQPasswordFromEnv
	}

	return nil
}

func (c *ConfigData) GetConfig() *Config {
	config := &Config{}
	if c.Server.Mode == "debug" {
		config.Port = c.Server.Debug.Port
		config.ReadTimeout = c.Server.Debug.ReadTimeout
		config.WriteTimeout = c.Server.Debug.WriteTimeout
		config.HandlerTimeout = c.Server.Debug.HandlerTimeout
	} else {
		config.Port = c.Server.Production.Port
		config.ReadTimeout = c.Server.Production.ReadTimeout
		config.WriteTimeout = c.Server.Production.WriteTimeout
		config.HandlerTimeout = c.Server.Production.HandlerTimeout
	}
	return config
}
