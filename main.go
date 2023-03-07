package main

import (
	"./slacknotif"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Handler_task(w http.ResponseWriter, r *http.Request) {
	user_wallet_address := "0x1234567890123456789012345678901234567890"

	task_name := slacknotif.TriggerNotification(user_wallet_address)

	fmt.Printf("Create Task : %s \n ", task_name)
}

func acknowledgeTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slacknotif.AcknowledgeTask(ctx, reqBody)

	w.WriteHeader(http.StatusOK)

}

func main() {

	http.HandleFunc("/app/v1/in_house/task_trigger", Handler_task)
	http.HandleFunc("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	http.ListenAndServe(":8080", nil)
}
