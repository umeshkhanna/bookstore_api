package app

import (
	"bookstore_api/logger"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("About to start the application....")
	router.Run(":8080")
}
