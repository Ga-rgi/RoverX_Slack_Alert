package slacknotification

import (
	"encoding/json"
	"os"

	"github.com/Ga-rgi/RoverX_Slack_Alert/dao"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

func AcknowledgeTask(reqBody []byte) {

	var taskPayload map[string]string
	err := json.Unmarshal(reqBody, &taskPayload)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		return
	}

	userWallet, ok := taskPayload["user_wallet_address"]
	if !ok {
		log.Error().Msg("missing user wallet address in request body")
		return
	}

	//check if the address is whitelisted
	isWhitelisted := dao.CheckAddressWhitelisted(userWallet)

	//log the address if it's not whitelisted
	if !isWhitelisted {
		log.Warn().Str("address", userWallet).Msg("unwhitelisted address")
	}

	SendSlackNotification(userWallet, dao.GetPartnerCommunity(userWallet))
}

func SendSlackNotification(address string, community string) {
	
	//Add slack token here
	slackClient := slack.New("xoxb-4911047364807-4949853786480-wrB9iQp7HkUAwZ77Dj9B1XdV")
	//Truncate address
	truncatedAddress := address[:4] + "..." + address[len(address)-4:]

	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	// Create notification message
	var notificationMessage string
	if community != "" {
		notificationMessage = truncatedAddress + " from " + community + " has joined RoverX "
		logger.Info().Str("community", community).Msgf(notificationMessage, truncatedAddress, community)
	} else {
		notificationMessage = truncatedAddress + " has joined RoverX"
		logger.Info().Msgf(notificationMessage, truncatedAddress)
	}
	// Send the notification to the configured Slack channel (#new-users)
	_, _, err := slackClient.PostMessage("#general", slack.MsgOptionText(notificationMessage, false))

	if err != nil {
		logger.Error().Err(err).Msg("failed to send Slack notification")
	}
}
