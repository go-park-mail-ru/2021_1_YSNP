package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logrusLogger *logrus.Entry
}

func NewLogger() *Logger {
	setFormatter()
	setLevel()

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

func setFormatter() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})
}

func setLevel() {
	logrus.SetLevel(logrus.DebugLevel)
}
