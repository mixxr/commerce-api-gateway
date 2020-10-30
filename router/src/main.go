package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

	router := createRouter() //gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/services/v1")
	{
		v1.GET("/:owner/:service", func(c *gin.Context) {
			name := c.Param("owner")
			action := c.Param("service")
			message := "searching for " + name + " is " + action
			c.String(http.StatusOK, message)
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
