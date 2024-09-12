package models

import (
	"log"
	"os"
	"strconv"

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

	overrideFromEnv("MODE", &c.Server.Mode)
	overrideFromEnv("RABBITMQ_HOST", &c.RabbitMQ.Host)
	overrideFromEnv("RABBITMQ_USER", &c.RabbitMQ.User)
	overrideFromEnv("RABBITMQ_PASS", &c.RabbitMQ.Password)

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

func overrideFromEnv(envKey string, configKey any) {
	value := os.Getenv(envKey)
	if value != "" {
		switch v := configKey.(type) {
		case *string:
			*v = value
			log.Printf("INFO : Override config with Env %s successfull", envKey)
		case *int:
			if valueInt, err := strconv.Atoi(value); err != nil {
				log.Printf("Error : Given env value %s of key %s cannot be parsed to int and connot be assigned to %v", value, envKey, configKey)
			} else {
				*v = valueInt
				log.Printf("INFO : Override config with Env %s successfull", envKey)
			}
		default:
			log.Printf("Error : Unsupported config type %v", configKey)
		}

	}
}
