package main

import (
	"net/http"
	"time"
	"turtle/core/dbclient"
	"turtle/core/lgr"
	"turtle/core/serverKit"
	"turtle/netes/deployListener"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	lgr.SetColors(true)
	lgr.SetOutputFolder("../logs", "TurtleNetes", true)

	serverKit.LoadGinConfig()
	dbclient.InitMongoDb()

	lgr.Info("Starting server with config: %+v", serverKit.SERVER_CONFIG)
	lgr.Info("Server URL: %s", serverKit.SERVER_CONFIG.GetURL())

	// Create Gin r
	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./static", true)))

	deployListener.InitDeployListenerApi(r)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:           serverKit.SERVER_CONFIG.GetAddress(),
		Handler:        r,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server
	lgr.Ok("Server is running at %s", serverKit.SERVER_CONFIG.GetURL())

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lgr.ErrorStack("Failed to start server: %v", err)
	}
}
