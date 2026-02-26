package auth

import (
	"errors"
	"fmt"
	"net/http"
	"turtle/core/lgr"
	serverKit2 "turtle/core/serverKit"

	"github.com/gin-gonic/gin"
)

func ApiKeysRequired(c *gin.Context) {

	header := c.GetHeader("Api-Key")

	user, hasUser := serverKit2.SERVER_CONFIG.ApiKeys[header]

	if hasUser {
		c.Set("userUid", user)
	} else {

		clientIP := c.ClientIP()

		// Bypass API key check for localhost
		if serverKit2.IsLocalhost(clientIP) {
			lgr.Info("Bypassing API key check for localhost request from: %s", clientIP)
			c.Set("userUid", "localhost")
			return
		} else {
			lgr.Error("Invalid api key: %s", header)
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid Api-Key"))
			return
		}

	}

}

func GetUserUidFromContext(c *gin.Context) (string, error) {
	user, userExist := c.Get("userUid")

	if userExist {
		return user.(string), nil
	} else {
		return "", errors.New("User Not Found")
	}

}
