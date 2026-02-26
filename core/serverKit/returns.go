package serverKit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnUnauthorized(c *gin.Context, err error) {
	c.String(http.StatusUnauthorized, err.Error())
}

func ReturnOkJson(c *gin.Context, jObj any) {
	c.JSON(http.StatusOK, jObj)
}

func ReturnUnacceptable(c *gin.Context, err error) {
	c.String(http.StatusNotAcceptable, err.Error())
}
func ReturnError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}
