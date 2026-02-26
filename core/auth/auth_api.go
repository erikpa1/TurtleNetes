package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
POST /api/auth/login
Body:

	{
	  "email": "user@mail.com",
	  "password": "secret"
	}
*/
func _TryToLoginUser(c *gin.Context) {
	// You already implemented this perfectly
	LoginHandler(c)
}

/*
POST /api/auth/activate
Body:

	{
	  "token": "activation-token"
	}
*/
func _TryToActivateUser(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO: implement activation logic
	// Example:
	// user, err := users.ActivateByToken(c.Request.Context(), req.Token)
	// if err != nil {
	//     c.AbortWithStatus(http.StatusUnauthorized)
	//     return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status": "account activated",
	})
}

func InitAuthApi(r *gin.Engine) {
	r.POST("/api/auth/login", _TryToLoginUser)
	r.POST("/api/auth/activate", _TryToActivateUser)
}
