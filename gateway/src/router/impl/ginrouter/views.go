package ginrouter

import (
	"net/http"

	"main/logger"

	"github.com/gin-gonic/gin"
)

func (o *GinRouter) servicesViewHandler(c *gin.Context) {

	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "/view/services"))

	c.HTML(http.StatusOK, "services", gin.H{ // services.html
		"apptitle": "Commerce Data Gateway", // TODO: from yaml

	})

	return
}

func (o *GinRouter) usersViewHandler(c *gin.Context) {

	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "/view/users"))

	c.HTML(http.StatusOK, "users", gin.H{ // users.html
		"apptitle": "Commerce Data Gateway", // TODO: from yaml

	})

	return
}
