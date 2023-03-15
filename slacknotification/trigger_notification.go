package slacknotification

import (
	"context"
	"encoding/json"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func TriggerNotification(userWallet string) {

	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msgf("Error loading .env file")
		return
	}

	projectID := os.Getenv("PROJECT_ID")
	queueID := os.Getenv("QUEUE_ID")
	locationID := os.Getenv("LOCATION_ID")
	url := os.Getenv("URL")

	createHTTPTask(projectID, queueID, locationID, url, userWallet)
}

func NewClient(ctx context.Context) *cloudtasks.Client {

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("NewClient: %v", err)
	}
	return client
}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(projectID, queueID, locationID, url, userWallet string) {

	ctx := context.Background()
	client := NewClient(ctx)
	if client == nil {
		log.Error().Err(err).Msgf("Failed to create client: %v", err)
		return
	}

	queuePath := "projects/" + projectID + "/locations/" + locationID + "/queues/" + queueID

	payload := map[string]string{
		"user_wallet_address": userWallet,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msgf("json.Marshal: %v", err)
		return

	}

	// Task payload.
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       payloadJSON,
				},
			},
		},
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		log.Error().Err(err).Msgf("client.CreateTask: %v", err)
		return
	}
	log.Info().Msgf("Task created with name %s", createdTask.Name)
	return
}
