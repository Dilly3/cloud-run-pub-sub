package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Environment string

const (
	ProdEnvironment    Environment = "prod"
	SandboxEnvironment Environment = "sandbox"
	DevEnvironment     Environment = "dev"
	PubSub             Environment = "pubsub"
)

type Configuration struct {
	Env  Environment `envconfig:"env" default:"dev"`
	Port string      `envconfig:"port" default:"8080"`
}

func GetConfig() (*Configuration, error) {
	var config Configuration
	err := envconfig.Process(string(PubSub), &config)
	if err != nil {
		return nil, err
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port = envPort
	}

	return &config, nil
}

func (c Configuration) IsProd() bool {
	return c.Env == ProdEnvironment
}

func (c Configuration) IsSandbox() bool {
	return c.Env == SandboxEnvironment
}

func (c Configuration) IsDev() bool {
	return c.Env == DevEnvironment
}
