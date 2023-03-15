package main

import (
	"Slack_notifs/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/Ga-rgi/RoverX_Slack_Alert/slacknotification"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func HandlerTask(c *gin.Context) {
	rate := limiter.Rate{
		Period: time.Minute,
		Limit:  10,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	limiterKey := c.ClientIP()
	context, err := limiter.Get(c.Request.Context(), limiterKey)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if context.Reached {
		c.Header("Retry-After", strconv.FormatInt(context.Reset, 10))
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	walletAddress := utils.HandleRequest(c)

	if walletAddress == "" || !utils.IsAddressValid(walletAddress) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	slacknotification.TriggerNotification(walletAddress)
	c.Status(http.StatusOK)
}

func acknowledgeTaskHandler(c *gin.Context) {
	walletAddress := utils.HandleRequest(c)
	slacknotification.AcknowledgeTask(walletAddress)
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()
	router.POST("/app/v1/in_house/task_trigger", HandlerTask)
	router.POST("/app/v1/in_house/acknowledge_task", acknowledgeTaskHandler)
	router.Run(":8080")
}
