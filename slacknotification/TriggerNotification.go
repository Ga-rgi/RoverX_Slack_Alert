package slacknotification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func TriggerNotification(userWallet string) string {
	projectID := "12345"
	queueID := "my-queue"
	locationID := "us-central1"
	url := "app/v1/in_house/acknowledge_task"

	task, err := createHTTPTask(projectID, locationID, queueID, url, userWallet)
	if err != nil {
		log.Fatalf("createHTTPTask: %v", err)
	}

	return task.GetName()

}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(projectID, locationID, queueID, url, userWallet string) (*taskspb.Task, error) {

	//Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	// Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	payload := map[string]string{
		"user_wallet_address": userWallet,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %v", err)
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
		return "", fmt.Errorf("client.CreateTask: %v", err)
	}

	return createdTask, nil
}
