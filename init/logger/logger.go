package logger

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var logDir string

type Logger struct {
	log *logrus.Logger
}

type Logging interface {
	Info(message string)
	InfoF(format string, args ...interface{})
	Debug(message string)
	DebugF(format string, args ...interface{})
	Error(message string)
	ErrorF(format string, args ...interface{})
	Fatal(message string)
	FatalF(format string, args ...interface{})
	Panic(message string)
	PanicF(format string, args ...interface{})
}

func newLoggerFile() (*os.File, error) {
	logFileName := time.Now().Format("2006-01-02") + ".log"
	p := path.Join(logDir, logFileName)
	logFile, err := os.OpenFile(p, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

func NewLogger(ctx context.Context, debug bool, dir string) (*Logger, error) {
	newLog := logrus.New()
	logDir = dir

	logFile, err := newLoggerFile()
	if err != nil {
		return nil, err
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	newLog.SetOutput(multiWriter)

	if debug {
		newLog.SetLevel(logrus.DebugLevel)
	}

	newLog.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})
	newLog.SetReportCaller(true)

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				newLogFile, err := newLoggerFile()
				if err != nil {
					continue
				}
				logFile.Close()
				newLog.SetOutput(newLogFile)
				logFile = newLogFile
			}
		}
	}()

	return &Logger{newLog}, nil
}

func (l *Logger) Info(message string) {
	l.log.Info(message)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Debug(message string) {
	l.log.Debug(message)
}

func (l *Logger) DebugF(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *Logger) Error(message string) {
	l.log.Error(message)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Fatal(message string) {
	l.log.Fatal(message)
}

func (l *Logger) FatalF(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *Logger) Panic(message string) {
	l.log.Panic(message)
}

func (l *Logger) PanicF(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
