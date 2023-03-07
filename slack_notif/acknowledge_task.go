package slacknotif

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/slack-go/slack"

	"github.com/Club-Defy/roverx-nft-backend/db/dao"
)

func AcknowledgeTask(ctx context.Context, reqBody []byte) {
	var taskPayload map[string]string
	err := json.Unmarshal(reqBody, &taskPayload)
	if err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return
	}

	userWallet, ok := taskPayload["user_wallet_address"]
	if !ok {
		log.Printf("Missing user wallet address in request body")
		return
	}

	//check if the address is whitelisted
	isWhitelisted := dao.CheckAddressWhitelisted(ctx, userWallet)

	//log the address if it's not whitelisted
	if !isWhitelisted {
		log.Printf("Unwhitelisted address: %s", userWallet)
	}

	//partner community for the user's address
	partnerCommunity := dao.GetPartnerCommunity(ctx, userWallet)
	slackClient := slack.New("SLACK_TOKEN")
	SendSlackNotification(ctx, slackClient, userWallet, partnerCommunity)
}

func SendSlackNotification(ctx context.Context, slackClient *slack.Client, address string, community string) error {
	// Truncate address
	truncatedAddress := fmt.Sprintf("%s...%s", address[:4], address[len(address)-4:])

	var notificationMessage string
	if community != "" {
		notificationMessage = fmt.Sprintf("%s from %s has joined RoverX", truncatedAddress, community)
	} else {
		notificationMessage = fmt.Sprintf("%s has joined RoverX", truncatedAddress)
	}

	// Send the notification to the configured Slack channel (#new-users)
	_, _, err := slackClient.PostMessageContext(ctx, "#new-users", slack.MsgOptionText(notificationMessage, false))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %v", err)
	}

	return nil
}
