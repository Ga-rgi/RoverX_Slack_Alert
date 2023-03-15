package slacknotification

import (
	"os"

	"github.com/Ga-rgi/RoverX_Slack_Alert/dao"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

func AcknowledgeTask(userWallet string) {

	isWhitelisted := dao.CheckAddressWhitelisted(userWallet)

	//log the address if it's not whitelisted
	if !isWhitelisted {
		log.Warn().Str("address", userWallet).Msg("non-whitelisted address")
	}
	SendSlackNotification(userWallet, dao.GetPartnerCommunity(userWallet))
}

func SendSlackNotification(address string, community string) {

	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msgf("Error loading .env file")
		return
	}

	slackClient := slack.New(os.Getenv("SLACK_TOKEN"))

	truncatedAddress := address[:4] + "..." + address[len(address)-4:]
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	var notificationMessage string
	if community != "" {
		notificationMessage = truncatedAddress + " from " + community + " has joined RoverX "
	} else {
		notificationMessage = truncatedAddress + " has joined RoverX"
	}

	// Send the notification to the configured Slack channel (#general)
	_, _, err = slackClient.PostMessage("#general", slack.MsgOptionText(notificationMessage, false))

	if err != nil {
		logger.Error().Err(err).Msg("failed to send Slack notification")
	}
}
