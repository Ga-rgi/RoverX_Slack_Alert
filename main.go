package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Ga-rgi/RoverX_Slack_Alert/slacknotification"
)

func Handler_task(w http.ResponseWriter, r *http.Request) {
	user_wallet_address := "0x4A906262CFE6B4de05A3E0b890Bf8eb4a4c2f30A"

	slacknotification.TriggerNotification(user_wallet_address)

}

func acknowledgeTaskHandler(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	slacknotification.AcknowledgeTask(reqBody)

	w.WriteHeader(http.StatusOK)

}

func main() {

	http.HandleFunc("/app/v1/in_house/task_trigger", Handler_task)
	http.HandleFunc("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	http.ListenAndServe(":8080", nil)
}
