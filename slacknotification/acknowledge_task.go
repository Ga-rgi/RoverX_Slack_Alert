package slacknotification

import (
	"encoding/json"
	"fmt"

	"github.com/Ga-rgi/RoverX_Slack_Alert/dao"
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

	//partner community for the user's address
	partnerCommunity := dao.GetPartnerCommunity(userWallet)

	SendSlackNotification(userWallet, partnerCommunity)
}

func SendSlackNotification(address string, community string) {

	slackClient := slack.New("xoxb-4911047364807-4949853786480-cu4JDuklOOajo4vqlCMYqSuv")
	//Truncate address
	truncatedAddress := address[:4] + "..." + address[len(address)-4:]

	var message string
	if community != "" {
		message = fmt.Sprintf("%s from %s has joined RoverX", truncatedAddress, community)
	} else {
		message = fmt.Sprintf("%s has joined RoverX", truncatedAddress)
	}

	// Send the notification to the configured Slack channel (#general)
	_, _, err := slackClient.PostMessage("#general", slack.MsgOptionText(message, false))
	if err != nil {
		log.Error().Err(err).Msg("failed to send Slack notification")
	}
}
