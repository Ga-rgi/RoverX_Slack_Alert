package main

import (
	"io/ioutil"
	"net/http"

	"github.com/Ga-rgi/RoverX_Slack_Alert/slacknotification"
	"github.com/gin-gonic/gin"
)

func Handler_task(c *gin.Context) {
	user_wallet_address := "0x4A906262CFE6B4de05A3E0b890Bf8eb4a4c2f30A"
	slacknotification.TriggerNotification(user_wallet_address)
	c.Status(http.StatusOK)
}

func acknowledgeTaskHandler(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	slacknotification.AcknowledgeTask(reqBody)
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()
	router.POST("/app/v1/in_house/task_trigger", Handler_task)
	router.POST("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	router.Run(":8080")
}
