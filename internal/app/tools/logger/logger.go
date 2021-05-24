package logger

import (
	formatters "github.com/fabienm/go-logrus-formatters"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrusLogger *logrus.Entry
}

func NewLogger(mode string, host string) *Logger {
	switch mode {
	case "prod":
		setProductionFormatter(host)
		setDevelopmentLevel()

	case "production":
		setProductionFormatter(host)
		setProductionLevel()

	case "development":
		setDevelopmentFormatter()
		setDevelopmentLevel()
	}

	return &Logger{
		logrusLogger: logrus.WithFields(logrus.Fields{
			"logger": "LOGRUS",
		}),
	}
}

func (l *Logger) GetLogger() *logrus.Entry {
	return l.logrusLogger
}

func GetDefaultLogger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	})
}

func (l *Logger) StartServerLog(host string, port string) {
	l.logrusLogger.WithFields(logrus.Fields{
		"host": host,
		"port": port,
	}).Info("Starting server")
}

func setProductionFormatter(host string) {
	logrus.SetFormatter(formatters.NewGelf(host))
}

func setDevelopmentFormatter() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})
}

func setProductionLevel() {
	logrus.SetLevel(logrus.InfoLevel)
}

func setDevelopmentLevel() {
	logrus.SetLevel(logrus.DebugLevel)
}

func (l *Logger) LogWSError(addr string, userID uint64, msg string) {
	l.logrusLogger.WithFields(logrus.Fields{
		"type":   "websocket",
		"client": addr,
		"user":   userID,
	}).Error(msg)
}

func (l *Logger) LogWSInfo(addr string, userID uint64, msg string) {
	l.logrusLogger.WithFields(logrus.Fields{
		"type":   "websocket",
		"client": addr,
		"user":   userID,
	}).Info(msg)
}
