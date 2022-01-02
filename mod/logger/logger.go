/**
*用于实现自定义logger
 */

package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type FileLogger struct {
	logger *logrus.Logger
}

const (
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	FatalLevel = logrus.FatalLevel
)

func New(logPath string, minLevel logrus.Level) (FileLogger, error) {

	file, err := os.Create(logPath)
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(minLevel)
	writers := []io.Writer{
		file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	log.SetOutput(fileAndStdoutWriter)
	if err != nil {
		return FileLogger{}, err
	}
	return FileLogger{
		logger: log,
	}, nil
}
func (f FileLogger) Debug(v ...interface{}) {
	f.logger.Debug(v...)
}
func (f FileLogger) Warn(v ...interface{}) {
	f.logger.Warn(v...)
}
func (f FileLogger) Error(v ...interface{}) {
	f.logger.Error(v...)
}
func (f FileLogger) Info(v ...interface{}) {
	f.logger.Info(v...)
}
func (f FileLogger) Debugf(format string, v ...interface{}) {
	f.logger.Debugf(fmt.Sprintf(format, v...))
}
func (f FileLogger) Warnf(format string, v ...interface{}) {
	f.logger.Warnf(fmt.Sprintf(format, v...))
}
func (f FileLogger) Errorf(format string, v ...interface{}) {
	f.logger.Errorf(fmt.Sprintf(format, v...))
}
func (f FileLogger) Infof(format string, v ...interface{}) {
	f.logger.Infof(fmt.Sprintf(format, v...))
}
func (f FileLogger) Sync() error {
	return f.Sync()
}
func output(v ...interface{}) string {
	_, file, line, _ := runtime.Caller(3)
	files := strings.Split(file, "/")
	file = files[len(files)-1]

	logFormat := "%s %s:%d " + fmt.Sprint(v...) + "\n"
	date := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(logFormat, date, file, line)
}

type Fields logrus.Fields

func (f FileLogger) WithFields(v logrus.Fields) *logrus.Entry {
	return f.logger.WithFields(v)
}
