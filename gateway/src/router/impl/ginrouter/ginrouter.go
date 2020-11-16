package ginrouter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"main/dataaccess"
	"main/dataaccess/models"
	sanitizer "main/router/utils/sanitizer"

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

func formatAndReturn(tables []models.ITable, err error, c *gin.Context, format string) {
	if err == nil {
		switch format {
		case "json":
			c.JSON(http.StatusOK, gin.H{
				"message": "TBD",
			})
		case "csv":
			var buffer bytes.Buffer

			for _, t := range tables {
				fmt.Fprintf(&buffer, "%s\n", t.String())
			}
			c.String(http.StatusOK, buffer.String())
		}
	} else {
		// ERROR
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
}

func toInt(s string, def int) int {
	intVar, err := strconv.Atoi(s)
	if err == nil {
		return intVar
	}
	return def
}

func toInt64(s string, def int64) int64 {
	intVar, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return intVar
	}
	return def
}

func getExt(c *gin.Context, paramname string, defval string) (string, string) {
	value := c.Param(paramname)
	i := strings.LastIndex(value, ".")
	if i < 0 {
		return value, defval
	} else {
		return value[:i], value[i+1:]
	}
}

func getInt(c *gin.Context, paramname string, defval int) int {
	value := c.Param(paramname)
	return toInt(value, defval)
}

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

	log.Println("Authentication:", username, password, hasAuth)

	c.Set(contextUsername, username)

	// Continue down the chain to handler etc
	c.Next()
}

func getStatus(c *gin.Context, owner string) int {
	username := c.GetString(contextUsername)
	log.Println("getStatus:", username)
	if username == owner {
		return models.StatusDraft
	}
	return models.StatusEnabled
}

func isOwner(c *gin.Context, owner string) bool {
	username := c.GetString(contextUsername)
	log.Println("isOwner:", username)

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
	router := createGinRouter()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/services/v1")
	//tables := v1.Group("/:owner/:service")
	v1.Use(AuthRequired)
	v1.Use(CreateLogEntry)
	{
		// eg. /abc/123/search/an_interesting_service/tag1/tag2/tag3
		// eg. /abc-/123-/search/an_interesting_service/tag1/tag2/tag3
		v1.GET("/:owner/:service/search/:descr/*tags", func(c *gin.Context) {
			owner := sanitizer.GetSearchToken(c.Param("owner"))     // can ends with -
			service := sanitizer.GetSearchToken(c.Param("service")) // can ends with -
			descr := sanitizer.GetDescr(c.Param("descr"))
			tags, ext := getExt(c, "tags", "csv")
			tags = sanitizer.GetTags(tags)
			//if owner != "" || service != "" || owner != "" || service != "" {
			tin := models.Table{Name: service, Owner: owner, Descr: descr, Tags: tags, Status: getStatus(c, owner)}
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET search", tin))
			if tin.IsValid() {
				tables, err := dal.ReadTables(&tin)
				//logger.CommonLog.Println("Result n.tables=", len(tables))
				formatAndReturn(models.ConvertToITables(tables), err, c, ext)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "owner and service have to be alphanumeric, use _ in place of spaces, use - in place of * for searches"})
				return
			}
		})
		v1.GET("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			if owner, service = sanitizer.CheckTokens(owner, service); owner != "" && service != "" {
				tin := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
				logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET service", tin))
				table, err := dal.ReadTable(&tin)
				formatAndReturn([]models.ITable{table}, err, c, ext)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "owner and service have to be alphanumeric, use _ in place of spaces"})
				return
			}
		})
		v1.GET("/:owner/:service/colnames/:lang", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang, ext := getExt(c, "lang", "csv")
			t := models.Table{Name: service, Owner: owner, DefLang: lang, Status: getStatus(c, owner)}
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicecolnames", t))
			table, err := dal.ReadTableColnames(&t, lang)
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicecolnames RESULT", table))
			formatAndReturn([]models.ITable{table}, err, c, ext)
		})
		v1.GET("/:owner/:service/values/:start/:count", func(c *gin.Context) {
			var startNum int
			var countNum int64
			owner := c.Param("owner")
			service := c.Param("service")
			startNum = getInt(c, "start", 0)
			count, ext := getExt(c, "count", "csv")
			countNum = toInt64(count, int64(dataAccessValueMaxCount))
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicevalues", t, startNum, countNum))
			table, err := dal.ReadTableValues(&t, startNum, countNum)
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "GET servicevalues RESULT", table))
			formatAndReturn([]models.ITable{table}, err, c, ext)
		})
		// CREATE
		v1.POST("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE service", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			var tableJson models.Table
			if err := c.ShouldBindJSON(&tableJson); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if tableJson.Name != service || tableJson.Owner != owner {
				c.JSON(http.StatusBadRequest, gin.H{"status": "service or owner does not match the url"})
				return
			}
			err := dal.StoreTable(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		v1.PUT("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "UPDATE service", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			var tableJson models.Table
			if err := c.ShouldBindJSON(&tableJson); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if tableJson.Name != service || tableJson.Owner != owner {
				c.JSON(http.StatusBadRequest, gin.H{"status": "service or owner does not match the url"})
				return
			}
			err := dal.UpdateTable(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		v1.POST("/:owner/:service/colnames/:lang", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang, ext := getExt(c, "lang", "csv")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE servicecolnames", owner, service, lang))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			var tableJson models.TableColnames
			if err := c.ShouldBindJSON(&tableJson); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tableJson.Lang = lang
			tableJson.SetParent(&models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)})

			err := dal.StoreTableColnames(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		// RETURN: number of affected rows in json format
		v1.POST("/:owner/:service/values", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "CREATE servicevalues", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			var tableJson models.TableValues
			if err := c.ShouldBindJSON(&tableJson); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tableJson.SetParent(&models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)})

			err := dal.StoreTableValues(&tableJson)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"count": tableJson.Count})
			return
		})
		// DELETE
		v1.DELETE("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE service", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err := dal.DeleteTable(&t)
			formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
		})
		// langs = / => all colnames
		// langs = /it => only it colnames
		// langs = /it/en/es => it,en and es colnames
		v1.DELETE("/:owner/:service/colnames/*langs", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE servicecolnames", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			langs, ext := getExt(c, "langs", "csv")
			langs = strings.Trim(langs, "/")
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err := dal.DeleteTableColnames(&t, strings.Split(langs, "/"))
			formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
			return
		})
		// count = 0 => ALL rows
		// count > 0 => TOP rows
		// count < 0 => BOTTOM rows
		v1.DELETE("/:owner/:service/values/*count", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			logger.AppLogger.Write(GetLogRequest(c, logger.LogInfo, "DELETE servicevalues", owner, service))
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			//startNum := getInt(c, "start", 0)
			count, ext := getExt(c, "count", "csv")
			countNum := toInt64(count, 0)
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err := dal.DeleteTableValues(&t, countNum)
			formatAndReturn([]models.ITable{&t}, err, c, ext)

			return
		})

	}

	return &GinRouter{
		dal:        dal,
		httpengine: router,
	}, nil
}

func (o *GinRouter) GetHandler() http.Handler {
	return o.httpengine
}
