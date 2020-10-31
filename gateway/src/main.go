package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func main() {
	fmt.Println("Starting Main===>\n.")
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	// router.Static("/assets", "./assets")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// data access
	var dal dataaccess.IDatastore
	var errSQL error
	dal, errSQL = dataaccess.NewDatastore()
	if errSQL != nil {
		fmt.Println("Starting Main===>FATAL ERROR:", errSQL)
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

			table, _ := dal.ReadTable(service, owner)
			c.String(http.StatusOK, table.String())
		})
		v1.GET("/:owner/:service/:start", func(c *gin.Context) {
			name := c.Param("owner")
			action := c.Param("service")
			start := c.Param("start")
			message := name + " is " + action + "," + start + ", -1"
			c.String(http.StatusOK, message)
		})
		v1.GET("/:owner/:service/:start/:count", func(c *gin.Context) {
			name := c.Param("owner")
			action := c.Param("service")
			start := c.Param("start")
			count := c.Param("count")
			message := name + " is " + action + "," + start + "," + count
			c.String(http.StatusOK, message)
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
	//fmt.Println("====> GO FUNC")
	if err := srv.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("graceful service shutdown")
		} else {
			fmt.Println("service cannot be started: ", err)
		}
	}
	//}()

	fmt.Println("<==== END")
}
