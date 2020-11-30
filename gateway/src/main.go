package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"main/dataaccess"
	"main/logger"
	"main/router"
	"main/router/impl/ginrouter"

	"github.com/spf13/viper"
)

type configuration struct {
	Runmode   string
	IpAndPort string
	Cert      string
	Key       string
}

func (o *configuration) String() string {
	var buffer bytes.Buffer
	val := reflect.ValueOf(o).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		buffer.WriteString(fmt.Sprintf("%s:%s, ", typeField.Name, valueField.Interface()))
	}

	return buffer.String()
}

var config *configuration

func loadConfigurations() {
	runmode, ok := os.LookupEnv("DCGW_RUNMODE")
	runmodeWasProvided := false
	if !ok {
		runmode = "dev"
	} else {
		runmodeWasProvided = true
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
		fmt.Println("Using default APP configurations...Error reading config file ", err)
		if runmodeWasProvided {
			fmt.Println("DCGW_RUNMODE was provided but not found as configuration file, provided runmode:", runmode)
			os.Exit(3)
		}
		config = &configuration{
			IpAndPort: ":8443",
			Cert:      "certs/server.crt",
			Key:       "certs/server.key",
		}
	} else {
		fmt.Println("Using provided APP configurations...")
		config = &configuration{
			IpAndPort: fmt.Sprintf("%s:%s",
				viper.Get("APP.IP").(string),
				viper.Get("APP.PORT").(string)),
			Cert: viper.Get("APP.CERT").(string),
			Key:  viper.Get("APP.KEY").(string),
		}
	}
	config.Runmode = runmode
}

func main() {

	logger.AppLogger.Info("main", "", "Starting Main ===>", logger.AppLogger)
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	logger.AppLogger.Info("main", "", mydir)

	loadConfigurations()
	logger.AppLogger.Info("main", "main", "configuration:", config)

	// data access
	var dal dataaccess.IDatastore
	var errSQL error
	dal, errSQL = dataaccess.NewDatastore()
	if errSQL != nil {
		reqLog := logger.LogRequest{
			UUID:     "main",
			Username: "main",
			ExitCode: 1,
			Message:  fmt.Sprintf("Dataaccess Error: %s", errSQL),
		}
		logger.AppLogger.WriteAndExit(&reqLog)
	}
	logger.AppLogger.Info("main", "main", "Data access layer: OK")

	// router
	var router router.IRouter
	var errRouter error
	router, errRouter = ginrouter.CreateRouter(dal)

	if errRouter != nil {
		reqLog := logger.LogRequest{
			UUID:     "main",
			Username: "main",
			ExitCode: 2,
			Message:  errRouter.Error(),
		}
		logger.AppLogger.WriteAndExit(&reqLog)
	}

	srv := &http.Server{
		Addr:           config.IpAndPort,
		Handler:        router.GetHandler(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.AppLogger.Info("main", "main", "ready to listen...", config.IpAndPort)
	fmt.Println("ready to listen...", config.IpAndPort)
	if err := srv.ListenAndServeTLS(config.Cert, config.Key); err != nil {
		if err == http.ErrServerClosed {
			logger.AppLogger.Info("main", "main", "service shutdown...")
		} else {
			reqLog := logger.LogRequest{
				UUID:     "main",
				Username: "main",
				ExitCode: 1,
				Message:  err.Error(),
			}
			logger.AppLogger.WriteAndExit(&reqLog)
		}
	}

	logger.AppLogger.Info("main", "main", "===> END.")
}
