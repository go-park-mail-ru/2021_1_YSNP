package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrusLogger *logrus.Entry
}

func NewLogger(mode string) *Logger {
	switch mode {
	case "production":
		setProductionFormatter()
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

func (l *Logger) StartServerLog(host string, port string) {
	l.logrusLogger.WithFields(logrus.Fields{
		"host": host,
		"port": port,
	}).Info("Starting server")
}

func setProductionFormatter() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
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