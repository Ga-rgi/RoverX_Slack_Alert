package slacknotification

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func TriggerNotification(userWallet string) {
	projectID := "12345"
	queueID := "my-queue"
	locationID := "us-central1"
	url := "app/v1/in_house/acknowledge_task"

	createHTTPTask(projectID, locationID, queueID, url, userWallet)
}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(projectID, locationID, queueID, url, userWallet string) {

	//Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("NewClient: %v", err)
		return
	}
	defer client.Close()

	// Task queue path.
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
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
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

}
