package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	ProdEnvironment    string = "prod"
	SandboxEnvironment string = "sandbox"
	DevEnvironment     string = "dev"
	PubSub             string = "pubs"
)

type Configuration struct {
	Env         string
	ProjectID   string
	TopicID     string
	Port        string
	PublishURL  string
	Location    string
	TaskQueueID string
	TaskDelay   int64
}

func GetConfig(logger *slog.Logger) (*Configuration, error) {
	var config Configuration
	err := godotenv.Load(".env_sample_var")
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
		return nil, err
	}

	envPort := os.Getenv("PORT")
	if envPort != "" {
		config.Port = envPort
	}

	envPublishURL := os.Getenv("PUBLISH_URL")
	if envPublishURL != "" {
		config.PublishURL = envPublishURL
	}

	envTaskQueueID := os.Getenv("TASK_QUEUE_ID")
	if envTaskQueueID != "" {
		config.TaskQueueID = envTaskQueueID
	}

	envLocation := os.Getenv("LOCATION")
	if envLocation != "" {
		config.Location = envLocation
	}

	envTaskDelay := os.Getenv("TASK_DELAY")
	if envTaskDelay != "" {
		config.TaskDelay, err = strconv.ParseInt(envTaskDelay, 10, 64)
		if err != nil {
			logger.Error("Error parsing TASK_DELAY", "error", err)
			return nil, err
		}
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
