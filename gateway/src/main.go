package main

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

	"dataaccess"
	"dataaccess/models"
	"router/utils"

	"github.com/gin-gonic/gin"
)

const (
	VALUES_MAX_COUNT int = 1000
)

func createRouter() *gin.Engine {
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
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

	return router
}

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

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {

	username, password, hasAuth := c.Request.BasicAuth()

	// TODO: check username/password
	// if hasAuth && username == "testuser" && password == "testpass" {
	// 	log.WithFields(log.Fields{
	// 		"user": username,
	// 	}).Info("User authenticated")
	// } else {
	// 	c.Abort()
	// 	c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	// 	return
	// }
	log.Println("Authentication:", username, password, hasAuth)

	c.Set("username", username)

	// Continue down the chain to handler etc
	c.Next()
}

func getStatus(c *gin.Context, owner string) int {
	username := c.GetString("username")
	log.Println("getStatus:", username)
	if username == owner {
		return models.StatusDraft
	}
	return models.StatusEnabled
}

func isOwner(c *gin.Context, owner string) bool {
	username := c.GetString("username")
	log.Println("isOwner:", username)

	return username == owner
}

var file *os.File

func main() {

	// log
	var err error
	file, err = os.OpenFile("main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Starting Main ===>")
	mydir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(mydir)

	// router.Static("/assets", "./assets")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// data access
	var dal dataaccess.IDatastore
	var errSQL error
	dal, errSQL = dataaccess.NewDatastore()
	if errSQL != nil {
		log.Println("Starting Main ===> FATAL ERROR:", errSQL)
		os.Exit(1)
	}
	log.Println("Data access layer: OK")

	// utils
	sanitizer := utils.NewSanitizer()

	// router
	router := createRouter() //gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/services/v1")
	//tables := v1.Group("/:owner/:service")
	v1.Use(AuthRequired)
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
			log.Println("Search:", tin)
			tables, errRead := dal.ReadTables(&tin)
			log.Println("Result n.tables=", len(tables))
			formatAndReturn(models.ConvertToITables(tables), errRead, c, ext)
		})
		v1.GET("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			if owner, service = sanitizer.CheckTokens(owner, service); owner != "" && service != "" {
				tin := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
				table, errRead := dal.ReadTable(&tin)
				formatAndReturn([]models.ITable{table}, errRead, c, ext)
			} else {
				formatAndReturn(nil, fmt.Errorf("owner and service names have to be alphanumeric, use _ in place of spaces"), c, "json")
			}

		})
		v1.GET("/:owner/:service/colnames/:lang", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang, ext := getExt(c, "lang", "csv")
			t := models.Table{Name: service, Owner: owner, DefLang: lang, Status: getStatus(c, owner)}
			log.Println("Table to search:", t)
			table, errRead := dal.ReadTableColnames(&t, lang)
			log.Println("Result:", table)
			formatAndReturn([]models.ITable{table}, errRead, c, ext)
		})
		v1.GET("/:owner/:service/values/:start/:count", func(c *gin.Context) {
			var startNum int
			var countNum int64
			owner := c.Param("owner")
			service := c.Param("service")
			startNum = getInt(c, "start", 0)
			count, ext := getExt(c, "count", "csv")
			countNum = toInt64(count, int64(VALUES_MAX_COUNT))
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			log.Println("Table to search:", t, startNum, countNum)
			table, errRead := dal.ReadTableValues(&t, startNum, countNum)
			log.Println("Result:", table)
			formatAndReturn([]models.ITable{table}, errRead, c, ext)
		})
		// CREATE
		v1.POST("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
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
			err = dal.StoreTable(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		v1.PUT("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
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
			err = dal.UpdateTable(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		v1.POST("/:owner/:service/colnames/:lang", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang, ext := getExt(c, "lang", "csv")
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

			err = dal.StoreTableColnames(&tableJson)
			formatAndReturn([]models.ITable{&tableJson}, err, c, ext)
		})
		// RETURN: number of affected rows in json format
		v1.POST("/:owner/:service/values", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
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

			err = dal.StoreTableValues(&tableJson)
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
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err = dal.DeleteTable(&t)
			formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
		})
		// langs = / => all colnames
		// langs = /it => only it colnames
		// langs = /it/en/es => it,en and es colnames
		v1.DELETE("/:owner/:service/colnames/*langs", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			langs, ext := getExt(c, "langs", "csv")
			langs = strings.Trim(langs, "/")
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err = dal.DeleteTableColnames(&t, strings.Split(langs, "/"))
			formatAndReturn([]models.ITable{&t}, err, c, ext) // TODO: StatusAccepted 202
			return
		})
		// count = 0 => ALL rows
		// count > 0 => TOP rows
		// count < 0 => BOTTOM rows
		v1.DELETE("/:owner/:service/values/*count", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			if !isOwner(c, owner) {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "you are not allowed"})
				return
			}
			//startNum := getInt(c, "start", 0)
			count, _ := getExt(c, "count", "csv")
			countNum := toInt64(count, 0)
			t := models.Table{Name: service, Owner: owner, Status: getStatus(c, owner)}
			err = dal.DeleteTableValues(&t, countNum)
			c.JSON(http.StatusAccepted, gin.H{"count": t.NRows})
			return
		})

	}

	srv := &http.Server{
		Addr:           ":8443",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//go func() {
	//log.Println("====> GO FUNC")
	if err := srv.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil {
		if err == http.ErrServerClosed {
			log.Println("graceful service shutdown")
		} else {
			log.Println("service cannot be started: ", err)
		}
	}
	//}()

	log.Println("<==== END")
}
