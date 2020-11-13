package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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
	file      *os.File
}

var AppLogger *Logger
var once sync.Once

func init() {
	once.Do(func() {
		// TODO: filename from env var
		filename := "application.log"
		openLogfile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("Error opening main log file:", err)
			os.Exit(1)
		}

		// TODO: level from env var
		AppLogger = &Logger{
			commonLog: log.New(openLogfile, "", log.Ldate|log.Ltime|log.Lmicroseconds),
			level:     LogInfo,
			file:      openLogfile,
			filename:  filename}

		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		AppLogger.commonLog.Println("Logger Init:", mydir, filename)
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
