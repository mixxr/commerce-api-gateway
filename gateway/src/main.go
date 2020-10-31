package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"dataaccess"
	"dataaccess/models"

	"github.com/gin-gonic/gin"
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

func formatAndReturnCSV(table models.ITable, err error, c *gin.Context) {
	if err == nil {
		c.String(http.StatusOK, table.String())
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
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

	// router
	router := createRouter() //gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/services/v1")
	{
		v1.GET("/:owner/:service", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			//var t models.Table{ Name: service, Owner:owner}
			table, errRead := dal.ReadTable(service, owner)
			formatAndReturnCSV(table, errRead, c)
		})
		v1.GET("/:owner/:service/colnames/:lang", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			lang := c.Param("lang")
			t := models.Table{Name: service, Owner: owner, DefLang: lang}
			log.Println("Table to search:", t)
			table, errRead := dal.ReadTableColnames(&t, lang)
			log.Println("Result:", table)
			formatAndReturnCSV(table, errRead, c)
		})
		v1.GET("/:owner/:service/values/:start/:count", func(c *gin.Context) {
			owner := c.Param("owner")
			service := c.Param("service")
			start := c.Param("start")
			count := c.Param("count")
			var countNum, startNum int
			startNum = toInt(start, 0)
			countNum = toInt(count, -1)
			t := models.Table{Name: service, Owner: owner}
			log.Println("Table to search:", t, startNum, countNum)
			table, errRead := dal.ReadTableValues(&t, startNum, countNum)
			log.Println("Result:", table)
			formatAndReturnCSV(table, errRead, c)
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
