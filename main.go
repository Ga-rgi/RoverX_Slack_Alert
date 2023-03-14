package main

import (
	"Slack_notifs/utils"
	"github.com/Ga-rgi/RoverX_Slack_Alert/slacknotification"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerTask(c *gin.Context) {
	walletAddress := utils.HandleRequest(c)
	if walletAddress != "" {
		slacknotification.TriggerNotification(walletAddress)
	}
	c.Status(http.StatusOK)
}

func acknowledgeTaskHandler(c *gin.Context) {
	walletAddress := utils.HandleRequest(c)
	if walletAddress != "" {
		slacknotification.AcknowledgeTask(walletAddress)
	}
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()
	router.POST("/app/v1/in_house/task_trigger", HandlerTask)
	router.POST("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	router.Run(":8000")
}
