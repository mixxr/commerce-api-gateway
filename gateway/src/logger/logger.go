package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/spf13/viper"
)

const (
	LogFatal   = 0
	LogError   = 1
	LogWarning = 2
	LogInfo    = 3
	LogDebug   = 4
	end        = 5
)

var logLevels = []string{
	"FATAL",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
}

type LogRequest struct {
	UUID     string
	Username string
	Level    int
	Message  string
	ExitCode int
	caller   string
	line     int
}

func (o *LogRequest) SetMessage(logmsgs ...interface{}) *LogRequest {
	o.Message = fmt.Sprint(logmsgs, " ")
	return o
}

func (o *LogRequest) SetLevel(l int) *LogRequest {
	o.Level = l
	return o
}

type Logger struct {
	commonLog *log.Logger
	level     int
	filename  string
	fileid    *os.File
}

var AppLogger *Logger
var once sync.Once

const (
	DefLogFileName = "application.log"
)

// func getGID() uint64 {
// 	b := make([]byte, 64)
// 	b = b[:runtime.Stack(b, false)]
// 	b = bytes.TrimPrefix(b, []byte("goroutine "))
// 	b = b[:bytes.IndexByte(b, ' ')]
// 	n, _ := strconv.ParseUint(string(b), 10, 64)
// 	return n
// }

func loadConfigurations() {
	runmode, ok := os.LookupEnv("DCGW_RUNMODE")
	if !ok {
		runmode = "dev"
	}

	// Set the fileid name of the configurations file
	viper.SetConfigName("config." + runmode + ".yaml")

	// Set the path to look for the configurations file
	configPath, ok := os.LookupEnv("DCGW_CONFIGPATH")
	if !ok {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(configPath)
	}

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Using default LOGGER configurations...Error reading config file, %s\n", err)
		AppLogger = &Logger{
			level:    LogInfo,
			filename: DefLogFileName}
	} else {
		fmt.Printf("Using env LOGGER configurations...%d\n", viper.Get("LOGGER.LEVEL").(int))
		viper.SetDefault("LOGGER.LEVEL", LogInfo)
		viper.SetDefault("LOGGER.FILENAME", DefLogFileName)
		viper.SetDefault("LOGGER.LOGPATH", ".")
		AppLogger = &Logger{
			level:    viper.Get("LOGGER.LEVEL").(int),
			filename: fmt.Sprintf("%s/%s", viper.Get("LOGGER.LOGPATH").(string), viper.Get("LOGGER.FILENAME").(string))}
	}
}

func init() {
	once.Do(func() {
		loadConfigurations()

		openLogfile, err := os.OpenFile(AppLogger.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Error opening main log fileid:", err)
			os.Exit(1)
		}

		AppLogger.commonLog = log.New(openLogfile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		AppLogger.fileid = openLogfile
		// TODO: defer AppLogger.fileid

		AppLogger.commonLog.Println("Logger Init:", AppLogger.filename, AppLogger.commonLog)
	})
}

func (o *Logger) GetLevel() int {
	return o.level
}

func (o *Logger) Write(logmsg *LogRequest) {
	if logmsg.Level < 0 || logmsg.Level >= end {
		logmsg.Level = LogInfo
	}
	if logmsg.Level <= o.level {
		if logmsg.caller == "" {
			_, logmsg.caller, logmsg.line, _ = runtime.Caller(1)
		}
		o.commonLog.Printf("[%s:%d] [%s] [%s] [%s] %s\n", filepath.Base(logmsg.caller), logmsg.line, logLevels[logmsg.Level], logmsg.Username, logmsg.UUID, logmsg.Message)
		logmsg.caller = ""
	}
}

func (o *Logger) WriteAndExit(logmsg *LogRequest) {
	logmsg.Level = 0
	_, logmsg.caller, logmsg.line, _ = runtime.Caller(1)
	o.Write(logmsg)
	os.Exit(logmsg.ExitCode)
}

func (o *Logger) Info(uuid string, username string, logmsgs ...interface{}) {
	_, caller, line, _ := runtime.Caller(1)
	o.Write(&LogRequest{caller: caller, line: line, UUID: uuid, Username: username, Level: LogInfo, Message: fmt.Sprint(logmsgs, " ")})
}

func (o *Logger) Fatal(uuid string, username string, logmsgs ...interface{}) {
	_, caller, line, _ := runtime.Caller(1)
	o.WriteAndExit(&LogRequest{caller: caller, line: line, UUID: uuid, Username: username, ExitCode: 2, Message: fmt.Sprint(logmsgs, " ")})
}
