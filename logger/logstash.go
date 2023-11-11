package logger

import (
	"fmt"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/softwareplace/go-logstash/env"
	"net"
	"os"
	"time"
)

func TimeInfoLogger() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02 15:04:05.99999 ")
}

var connection net.Conn
var currentFilePath string
var logger *logrus.Logger

func Logger(loggerName string) *logrus.Entry {
	var logPath = env.GetEnv(env.LoggerPath, "/var/log/"+env.GetAppName())

	newLogDir := fmt.Sprintf("%s/%s/", logPath, time.Now().Format("2006-01"))
	logFileName := fmt.Sprintf("%s%s-%s.log", newLogDir, env.GetAppName(), time.Now().Format("2006-01-02"))

	if logger != nil || (currentFilePath != "" && currentFilePath != logFileName) {
		return logger.WithFields(logrus.Fields{
			"date":             TimeInfoLogger(),
			"application_name": env.GetAppName(),
			"logger_name":      loggerName,
		})
	}

	currentFilePath = logFileName

	logger = logrus.New()

	rotateLogs, _ := rotatelogs.New(
		logFileName,
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	// set the new rotateLogs hook to the logger
	logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  rotateLogs,
			logrus.FatalLevel: rotateLogs,
		},
		&logrus.JSONFormatter{},
	))

	connection = connectionCreate(logger)

	if connection != nil {
		hook := logrustash.New(connection, &logrus.JSONFormatter{})
		logger.Hooks.Add(hook)
	}

	return logger.WithFields(logrus.Fields{
		"date":             TimeInfoLogger(),
		"application_name": env.GetAppName(),
		"logger_name":      loggerName,
	})
}

func connectionCreate(log *logrus.Logger) net.Conn {
	isLogstashEnable := env.GetEnvBool(env.LogstashEnable, false)

	if isLogstashEnable {
		logstashUri := os.Getenv(env.LogstashUri)

		if logstashUri != "" {
			timeout := time.Second * time.Duration(env.GetEnvAsInt(env.LogstashTimeout, 5))

			conn, err := net.DialTimeout("tcp", logstashUri, timeout)

			if err != nil {
				log.Error(TimeInfoLogger(), "Failed to connect to logstash", err)
			}
			return conn
		} else {
			log.Warn(TimeInfoLogger(), "Trying to create a logstash connection but LOGSTASH_URI was not found")
		}

	}
	return nil
}
