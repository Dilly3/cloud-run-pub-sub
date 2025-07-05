package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

const (
	ProdEnvironment    string = "prod"
	SandboxEnvironment string = "sandbox"
	DevEnvironment     string = "dev"
	PubSub             string = "pubs"
)

type Configuration struct {
	Env       string
	Port      string
	ProjectID string
	TopicID   string
}

func GetConfig(logger *slog.Logger) (*Configuration, error) {
	var config Configuration
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
		return nil, err
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port = envPort
	}
	envProjectID := os.Getenv("PROJECT_ID")
	if envProjectID != "" {
		config.ProjectID = envProjectID
	}
	envTopicID := os.Getenv("TOPIC_ID")
	if envTopicID != "" {
		config.TopicID = envTopicID
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
