package utils

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HandleRequest(c *gin.Context) string {

	var taskPayload map[string]string
	if err := c.BindJSON(&taskPayload); err != nil {
		log.Error().Err(err).Msg("failed to parse request body")
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Failed to read request body"})
		return ""
	}

	userWallet, ok := taskPayload["user_wallet_address"]
	if !ok {
		log.Error().Msg("missing user wallet address in request body")
		return ""
	}
	return userWallet
}

func IsAddressValid(userWallet string) bool {
	AddressRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return AddressRegex.MatchString(userWallet)
}
