package deployListener

import (
	"turtle/core/serverKit"
	"turtle/netes/netesAuth"

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
	r.GET("/deplistener/ping", netesAuth.NetesApiKeyRequired, _Ping)
	r.GET("/deplistener/info", netesAuth.NetesApiKeyRequired, _GetInfo)
	r.POST("/deplistener/info", netesAuth.NetesApiKeyRequired, _PostInfo)

	r.POST("/deplistener/receive", _ReceiveDeploymentPackage)
}
