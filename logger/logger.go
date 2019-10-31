package logger

import (
	"io"
	"os"

	"git.corp.adobe.com/dc/notifications_load_test/config"
	"git.corp.adobe.com/dc/notifications_load_test/constants"
	"gopkg.in/inconshreveable/log15.v2"
)

var rootLog = log15.New("service", constants.ProgramName)

func init() {
	logFmt := log15.LogfmtFormat()
	AddStackTraceLogging(rootLog, os.Stdout, logFmt, nil)
}

// New creates new instance of logger
func New(kvps ...interface{}) log15.Logger {
	return rootLog.New(kvps...)
}

// AddStackTraceLogging will add logging filters depending
// upon LogLevel specified in the config file
func AddStackTraceLogging(log log15.Logger, writer io.Writer, logFmt log15.Format,
	appCfg *config.AppConfig) {

	var cfg *config.AppConfig

	if appCfg == nil {
		cfg = config.GetConfig()
	} else {
		cfg = appCfg
	}

	// for logs with LogLevel warning to critical it will only display stack trace
	// for logs with LogLevel info and debug it will include both stack trace and info logs
	switch level := cfg.LogLevel; level {
	case "critical":
		stackHandler := getStackHandler(writer, logFmt, log15.LvlCrit)
		log.SetHandler(log15.MultiHandler(stackHandler))
	case "error":
		stackHandler := getStackHandler(writer, logFmt, log15.LvlError)
		log.SetHandler(log15.MultiHandler(stackHandler))
	case "warning":
		stackHandler := getStackHandler(writer, logFmt, log15.LvlWarn)
		log.SetHandler(log15.MultiHandler(stackHandler))
	default:
		stackHandler := getStackHandler(writer, logFmt, log15.LvlWarn)
		infoHandler := getInfoHandler(writer, logFmt, level)
		log.SetHandler(log15.MultiHandler(stackHandler, infoHandler))
	}
}

// getStackHandler method appends stacktrace to logs
func getStackHandler(writer io.Writer, logFmt log15.Format, lvl log15.Lvl) log15.Handler {
	stackHandler := log15.StreamHandler(writer, logFmt)
	stackHandler = log15.CallerStackHandler("%+v", stackHandler)
	stackHandler = log15.FilterHandler(func(r *log15.Record) bool {
		return r.Lvl <= lvl
	}, stackHandler)
	return stackHandler
}

// getInfoHandler method appends all logs if lvl is debug
// otherwise assume lvl is info and return info logs only
func getInfoHandler(writer io.Writer, logFmt log15.Format, level string) log15.Handler {
	infoHandler := log15.StreamHandler(writer, logFmt)
	if level == "debug" {
		infoHandler = log15.FilterHandler(func(r *log15.Record) bool {
			return r.Lvl >= log15.LvlInfo
		}, infoHandler)
	} else {
		infoHandler = log15.FilterHandler(func(r *log15.Record) bool {
			return r.Lvl == log15.LvlInfo
		}, infoHandler)
	}
	return infoHandler
}
