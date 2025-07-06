package publisher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	pubsb "cloud.google.com/go/pubsub"
	"github.com/dilly3/cloud-run-pub-sub/config"
	"github.com/dilly3/cloud-run-pub-sub/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Publisher struct {
	projectID string
	topicID   string
	logger    *slog.Logger
	config    *config.Configuration
}

func NewPublisher(projectID, topicID string, logger *slog.Logger, config *config.Configuration) *Publisher {
	return &Publisher{
		projectID: projectID,
		topicID:   topicID,
		logger:    logger,
		config:    config,
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

func (p *Publisher) QueueTask(data any, delayInSeconds int64) (string, error) {
	task, err := p.createHTTPTask(p.config.ProjectID, p.config.Location, p.config.TaskQueueID, p.config.PublishURL, delayInSeconds, data)
	if err != nil {
		p.logger.Error("Failed to create HTTP task", "error", err)
		return "", fmt.Errorf("createHTTPTask: %w", err)
	}

	return task.Name, nil
}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func (p *Publisher) createHTTPTask(projectID, locationID, taskQueueID, url string, delayInSeconds int64, message any) (*taskspb.Task, error) {

	// Create a new Cloud Tasks client instance.
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, taskQueueID)
	body, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	// Build the Task payload.
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Body:       body,
				},
			},
			ScheduleTime: timestamppb.New(time.Now().Add(time.Duration(delayInSeconds) * time.Second)),
		},
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}

	return createdTask, nil
}

func (p *Publisher) Publish(data any) (string, error) {

	ctx := context.Background()
	client, err := pubsb.NewClient(ctx, p.projectID)
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
	result := t.Publish(ctx, &pubsb.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("pubsub: result.Get: %w", err)
	}
	return id, nil
}

func (p *Publisher) DecodePubSubData(body io.ReadCloser, data any) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// parse as a Pub/Sub push message
	var pubSubMessage models.PubSubMessage
	if err := json.Unmarshal(bodyBytes, &pubSubMessage); err == nil {
		if pubSubMessage.Message.Data == "" {
			return fmt.Errorf("Pub/Sub message has no data")
		}

		// Decode the base64 data from Pub/Sub message
		decodedData, err := decodePubSubData(pubSubMessage.Message.Data)
		if err != nil {
			return fmt.Errorf("failed to decode base64 data from Pub/Sub message: %w", err)
		}

		if err := json.Unmarshal(decodedData, &data); err != nil {
			return fmt.Errorf("failed to unmarshal transaction from Pub/Sub data: %w", err)
		}
		return nil
	}
	return fmt.Errorf("invalid Pub/Sub message")
}

func (p *Publisher) DecodeQueueData(body io.ReadCloser, data any) error {
	// Read the request body
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	defer body.Close()

	// Unmarshal the JSON data
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// decodePubSubData decodes base64 data from Pub/Sub message
func decodePubSubData(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// PublishWithDelay publishes a message with an optional delay
// delay is the duration to wait before the message becomes available for delivery
// func (p *Publisher) publishWithDelay(data any, delayInSeconds int64) (string, error) {
// 	ctx := context.Background()
// 	client, err := pubsub.NewClient(ctx, p.projectID)
// 	if err != nil {
// 		return "", fmt.Errorf("pubsub: NewClient: %w", err)
// 	}
// 	defer client.Close()

// 	// Extract just the topic name from the full path
// 	topicName := extractTopicName(p.topicID)
// 	t := client.Topic(topicName)
// 	msg, err := json.Marshal(data)
// 	if err != nil {
// 		return "", fmt.Errorf("json: Marshal: %w", err)
// 	}

// 	// Create message with optional publish time for delay
// 	message := &pubsub.Message{
// 		Data: []byte(msg),
// 	}

// 	// If delay is specified, set the publish time to the future
// 	if delayInSeconds > 0 {
// 		publishTime := time.Now().Add(time.Duration(delayInSeconds) * time.Second)
// 		message.PublishTime = publishTime
// 	}

// 	result := t.Publish(ctx, message)
// 	// Block until the result is returned and a server-generated
// 	// ID is returned for the published message.
// 	id, err := result.Get(ctx)
// 	if err != nil {
// 		return "", fmt.Errorf("pubsub: result.Get: %w", err)
// 	}
// 	return id, nil
// }
