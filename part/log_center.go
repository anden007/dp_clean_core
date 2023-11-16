package part

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// "disable" < "fatal" < "error" < "warn" < "info" < "debug"
func parseLogLevel(level string) (result logrus.Level) {
	switch level {
	case "fatal":
		result = logrus.FatalLevel
	case "error":
		result = logrus.ErrorLevel
	case "warn":
		result = logrus.WarnLevel
	case "info":
		result = logrus.InfoLevel
	case "debug":
		result = logrus.DebugLevel
	}
	return
}

type ILogCenter interface {
	Log() *logrus.Logger
}

type LogCenter struct {
	logger *logrus.Logger
}

func NewLogCenter() ILogCenter {
	writers := []io.Writer{}
	logLevel := viper.GetString("log.level")
	serverName := os.Getenv("GODF_SERVER_NAME")
	if serverName == "" {
		serverName = "Unknow"
	}
	logFilePath := "./logs"
	if _, err := os.Stat(logFilePath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logFilePath, 0777)
		}
	}
	instance := &LogCenter{
		logger: logrus.New(),
	}
	instance.logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	instance.logger.Level = parseLogLevel(logLevel)
	rotateWriter, _ := rotatelogs.New(
		fmt.Sprintf("%s/%%Y%%m%%d.log", logFilePath),
		rotatelogs.WithRotationCount(0),
		rotatelogs.WithRotationTime(time.Hour*24), //24小时分割文件，没有找到按天零点切割的配置，如有需要可自行实现
	)
	writers = append(writers, rotateWriter)
	if ENV == ENUM_ENV_DEV {
		writers = append(writers, os.Stdout)
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	instance.logger.Out = fileAndStdoutWriter
	return instance
}

func (m *LogCenter) Log() *logrus.Logger {
	return m.logger
}
