package ginrouter

import (
	"fmt"
	"io"
	"os"

	"time"

	"main/dataaccess"
	"main/dataaccess/models"

	"main/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const (
	dataAccessValueMaxCount int    = 1000
	contextXRequestID       string = "X-Request-Id"
	contextUsername         string = "username"
	contextLogEntry         string = "logEntry"
)

func loadConfigurations() {
	runmode, ok := os.LookupEnv("DCGW_RUNMODE")
	if !ok {
		runmode = "dev"
	}
	logger.AppLogger.Info("main", "main", "runmode:", runmode)

	// Set the file name of the configurations file
	viper.SetConfigName("config." + runmode + ".yaml")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using default ROUTER configurations...Error reading config file ", err)
	} else {
		fmt.Println("Using provided ROUTER configurations...")
	}
}

func RequestIdMiddleware(c *gin.Context) {
	reqID := uuid.New().String()
	c.Writer.Header().Set(contextXRequestID, reqID)
	c.Set(contextXRequestID, reqID)
	c.Next()
}

// type LogRequest struct {
// 	UUID     string
// 	Username string
// 	Level    int
// 	Message  string
// 	ExitCode int
// }
func CreateLogEntry(c *gin.Context) {
	reqLog := logger.LogRequest{
		UUID:     c.GetString(contextXRequestID),
		Username: c.GetString(contextUsername),
	}
	c.Set(contextLogEntry, &reqLog)
	c.Next()
}

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {

	username, password, hasAuth := c.Request.BasicAuth()

	logger.AppLogger.Info(c.GetString(contextXRequestID), "main", "Authentication:", username, password, hasAuth)

	c.Set(contextUsername, username)

	// Continue down the chain to handler etc
	c.Next()
}

func getStatus(c *gin.Context, owner string) int {
	username := c.GetString(contextUsername)
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "getStatus:", username))
	if username == owner {
		return models.StatusDraft
	}
	return models.StatusEnabled
}

func isOwner(c *gin.Context, owner string) bool {
	username := c.GetString(contextUsername)
	logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "isOwner:", username))

	return username == owner
}

func GetLogRequest(c *gin.Context, level int, logmsgs ...interface{}) (lr *logger.LogRequest) {
	lr = c.MustGet(contextLogEntry).(*logger.LogRequest)
	lr.Message = fmt.Sprint(logmsgs, " ")
	lr.Level = level
	return
}

//var file *os.File

func createGinRouter() *gin.Engine {
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("requests.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("%s - [%s] [%s] [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Keys[contextUsername],
			param.Keys[contextXRequestID],
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.Use(gin.Recovery()) // 500 management
	router.Use(RequestIdMiddleware)

	return router
}

type GinRouter struct {
	dal        dataaccess.IDatastore
	httpengine *gin.Engine
}

func CreateRouter(dal dataaccess.IDatastore) (*GinRouter, error) {

	if dal == nil {
		return nil, fmt.Errorf("data access is needed.")
	}
	loadConfigurations()

	// gin router
	//router := createGinRouter()
	ginRouter := &GinRouter{
		dal:        dal,
		httpengine: createGinRouter(),
	}

	ginRouter.httpengine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := ginRouter.httpengine.Group("/services/v1")
	//tables := v1.Group("/:owner/:service")
	v1.Use(AuthRequired)
	v1.Use(CreateLogEntry)
	{
		// eg. /abc/123/search/an_interesting_service/tag1/tag2/tag3
		// eg. /abc-/123-/search/an_interesting_service/tag1/tag2/tag3
		v1.GET("/:owner/:service/search/:descr/*tags", ginRouter.searchHandler)
		v1.GET("/:owner/:service", ginRouter.getServiceHandler)
		v1.GET("/:owner/:service/colnames/:lang", ginRouter.getColnamesHandler)
		v1.GET("/:owner/:service/values/:start/:count", ginRouter.getValuesHandler)
		// CREATE
		v1.POST("/:owner/:service", ginRouter.postServiceHandler)
		v1.POST("/:owner/:service/colnames/:lang", ginRouter.postColnamesHandler)
		v1.POST("/:owner/:service/values", ginRouter.postValuesHandler) // RETURN: number of affected rows in json format
		// UPDATE
		v1.PUT("/:owner/:service", ginRouter.putServiceHandler)
		// DELETE
		v1.DELETE("/:owner/:service", ginRouter.deleteServiceHandler)
		// langs = / => all colnames
		// langs = /it => only it colnames
		// langs = /it/en/es => it,en and es colnames
		v1.DELETE("/:owner/:service/colnames/*langs", ginRouter.deleteColnamesHandler)
		// count = 0 => ALL rows
		// count > 0 => TOP rows
		// count < 0 => BOTTOM rows
		v1.DELETE("/:owner/:service/values/*count", ginRouter.deleteValuesHandler)

	}

	return ginRouter, nil
}
