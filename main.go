
/* Two programs effectively execute the requirement to send Slack notification.
slacknotifs/TriggerNotification creates a task in the task queue 
slacknotifs/ acknowledge_task gets the user address, checks if its whitelisted and sends the slack alert */

package main

import (
	"./slacknotif"
	"fmt"
	"io/ioutil"
	"net/http"
)

//sends a user wallet address to TriggerNotification (performed by wallet_service as a part of the user sign up workflow)
func Handler_task(w http.ResponseWriter, r *http.Request) {
	user_wallet_address := "0x1234567890123456789012345678901234567890"
	
	task_name := slacknotif.TriggerNotification(user_wallet_address)

	fmt.Printf("Create Task : %s \n ", task_name)
}

// Gets user wallet from the Task queue
func acknowledgeTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
        //calls Acknowledgetask for the slack alert
	slacknotif.AcknowledgeTask(ctx, reqBody)

	w.WriteHeader(http.StatusOK)

}

func main() {

	http.HandleFunc("/app/v1/in_house/task_trigger", Handler_task) //handles route to TriggerNotification
	http.HandleFunc("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	http.ListenAndServe(":8080", nil)
}
