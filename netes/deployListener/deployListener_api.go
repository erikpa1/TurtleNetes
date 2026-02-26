package deployListener

import (
	"turtle/core/serverKit"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func _Ping(c *gin.Context) {
	serverKit.ReturnOkJson(c, bson.M{"status": "ok"})
}

func _GetInfo(c *gin.Context) {
	serverKit.ReturnOkJson(c, bson.M{"status": "ok"})
}

func _PostInfo(c *gin.Context) {
	serverKit.ReturnOkJson(c, bson.M{"status": "ok"})
}

func _ReceiveDeploymentPackage(c *gin.Context) {
	serverKit.ReturnOkJson(c, bson.M{"status": "ok"})
}

func InitDeployListenerApi(r *gin.Engine) {
	r.GET("/deplistener/ping", _Ping)
	r.GET("/deplistener/info", _GetInfo)
	r.POST("/deplistener/info", _PostInfo)

	r.POST("/deplistener/receive", _ReceiveDeploymentPackage)
}
