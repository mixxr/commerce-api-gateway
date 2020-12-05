package ginrouter

import (
	"net/http"

	"main/logger"

	"github.com/gin-gonic/gin"
)

func (o *GinRouter) indexViewHandler(c *gin.Context) {

	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "/view/services"))

	c.HTML(http.StatusOK, "services/index", gin.H{
		"apptitle": "Commerce Data Gateway", // TODO: from yaml

	})

	return
}
