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

func main() {

	// log
	file, err := os.OpenFile("main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Hello world!")

	log.Println("Starting Main===>\n.")
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
		log.Println("Starting Main===>FATAL ERROR:", errSQL)
		os.Exit(1)
	}

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
	{
		// eg. /abc/123/search/an_interesting_service/tag1/tag2/tag3
		// eg. /abc*/123*/search/an_interesting_service/tag1/tag2/tag3
		v1.GET("/:owner/:service/search/:descr/*tags", func(c *gin.Context) {
			owner := sanitizer.GetSearchToken(c.Param("owner"))     // can ends with -
			service := sanitizer.GetSearchToken(c.Param("service")) // can ends with -
			descr := sanitizer.GetDescr(c.Param("descr"))
			tags, ext := getExt(c, "tags", "csv")
			tags = sanitizer.GetTags(tags)
			//if owner != "" || service != "" || owner != "" || service != "" {
			tin := models.Table{Name: service, Owner: owner, Descr: descr, Tags: tags}
			log.Println("Search:", tin)
			tables, errRead := dal.ReadTables(&tin)
			log.Println("Result n.tables=", len(tables))
			formatAndReturn(models.ConvertToITables(tables), errRead, c, ext)
		})
		v1.GET("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service, ext := getExt(c, "service", "csv")
			if owner, service = sanitizer.CheckTokens(owner, service); owner != "" && service != "" {
				tin := models.Table{Name: service, Owner: owner}
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
			t := models.Table{Name: service, Owner: owner, DefLang: lang}
			log.Println("Table to search:", t)
			table, errRead := dal.ReadTableColnames(&t, lang)
			log.Println("Result:", table)
			formatAndReturn([]models.ITable{table}, errRead, c, ext)
		})
		v1.GET("/:owner/:service/values/:start/:count", func(c *gin.Context) {
			var countNum, startNum int
			owner := c.Param("owner")
			service := c.Param("service")
			startNum = getInt(c, "start", 0)
			count, ext := getExt(c, "count", "csv")
			countNum = toInt(count, VALUES_MAX_COUNT)
			t := models.Table{Name: service, Owner: owner}
			log.Println("Table to search:", t, startNum, countNum)
			table, errRead := dal.ReadTableValues(&t, startNum, countNum)
			log.Println("Result:", table)
			formatAndReturn([]models.ITable{table}, errRead, c, ext)
		})
		// CREATE
		v1.POST("/:owner/:service", func(c *gin.Context) {
			var tableJson models.Table
			if err := c.ShouldBindJSON(&tableJson); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if tableJson.Name == "" || tableJson.Owner == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "name or owner is empty"})
				return
			}
			message := "CREATED service " + tableJson.Name + ", by " + tableJson.Owner
			c.String(http.StatusOK, message)
		})
		v1.PUT("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			message := "MODIFIED service " + service + ", by " + owner
			c.String(http.StatusOK, message)
		})
		v1.POST("/:owner/:service/colnames", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			message := "ADDED colnames to " + service + ", by " + owner
			c.String(http.StatusOK, message)
		})
		v1.POST("/:owner/:service/values", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			message := "ADDED values to " + service + ", by " + owner
			c.String(http.StatusOK, message)
		})
		// DELETE
		v1.DELETE("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			message := "REMOVED service " + service + ", by " + owner
			c.String(http.StatusOK, message)
		})
		// langs = / => all colnames
		// langs = /it => only it colnames
		// langs = /it/en/es => it,en and es colnames
		v1.DELETE("/:owner/:service/colnames/*langs", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang := c.Param("langs")
			message := "REMOVED colnames " + lang + " for " + service + ", by " + owner
			c.String(http.StatusOK, message)
		})
		// start =0, count = -1 => ALL rows
		v1.DELETE("/:owner/:service/values/:start/:count", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			start := c.Param("start")
			count := c.Param("count")
			message := "REMOVED values from " + start + " count=" + count + " for " + service + ", by " + owner
			c.String(http.StatusOK, message)
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
