package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pubsub "cloud.google.com/go/pubsub"
)

type Publisher struct {
	projectID string
	topicID   string
	logger    *slog.Logger
}

func NewPublisher(projectID, topicID string, logger *slog.Logger) *Publisher {
	return &Publisher{
		projectID: projectID,
		topicID:   topicID,
		logger:    logger,
	}
}

// extractTopicName extracts just the topic name from a full topic path
// e.g., "projects/my-project/topics/my-topic" -> "my-topic"
func extractTopicName(topicPath string) string {
	if strings.Contains(topicPath, "/topics/") {
		parts := strings.Split(topicPath, "/topics/")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return topicPath
}

func (p *Publisher) Publish(data any, delayInSeconds int64) (string, error) {
	return p.publishWithDelay(data, delayInSeconds)
}

// PublishWithDelay publishes a message with an optional delay
// delay is the duration to wait before the message becomes available for delivery
func (p *Publisher) publishWithDelay(data any, delayInSeconds int64) (string, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, p.projectID)
	if err != nil {
		return "", fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	// Extract just the topic name from the full path
	topicName := extractTopicName(p.topicID)
	t := client.Topic(topicName)
	msg, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json: Marshal: %w", err)
	}

	// Create message with optional publish time for delay
	message := &pubsub.Message{
		Data: []byte(msg),
	}

	// If delay is specified, set the publish time to the future
	if delayInSeconds > 0 {
		publishTime := time.Now().Add(time.Duration(delayInSeconds) * time.Second)
		message.PublishTime = publishTime
	}

	result := t.Publish(ctx, message)
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("pubsub: result.Get: %w", err)
	}
	return id, nil
}
