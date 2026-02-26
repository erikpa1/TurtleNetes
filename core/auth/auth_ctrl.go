package auth

import (
	"net/http"
	"turtle/core/serverKit"
	"turtle/users"

	"github.com/gin-gonic/gin"
)

func LoginOrLocalhost(c *gin.Context) {
	clientIP := c.ClientIP()

	// 1️⃣ Localhost bypass
	if serverKit.IsLocalhost(clientIP) {
		user := users.User{
			Email: "localhost@pointe.sk",
			Role:  "superadmin",
		}

		c.Set("userUid", "localhost")
		c.Set("user", user)
		c.Next()
		return
	}

	// 2️⃣ JWT from cookie
	tokenStr, err := c.Cookie(JWT_COOKIE_NAME)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	claims, err := ParseToken(tokenStr, clientIP)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	user := users.User{
		Email: claims.Email,
		Role:  claims.Role,
	}

	c.Set("userUid", claims.UserID)
	c.Set("user", user)

	c.Next()
}
