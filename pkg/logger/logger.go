package logger

import (
	"fmt"
	"log"
	"time"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type defaultLogger struct{}

var DefaultLogger = defaultLogger{}

func InfoLabel() string {
	return WithColor(Cyan, "[INFO]")
}

func ErrorLabel() string {
	return WithColor(Red, "[ERROR]")
}

func WarnLabel() string {
	return WithColor(Yellow, "[WARN]")
}

func DateLabel() string {
	return DateFormat(time.Now())
}

func DateFormat(t time.Time) string {
	return WithColor(Gray, t.Format("2006-01-02T15:04:05 -0700"))
}

func (l defaultLogger) SInfof(format string, args ...any) string {
	return fmt.Sprintf(fmt.Sprintf("%s %s %s\n", InfoLabel(), DateLabel(), format), args...)
}

func (l defaultLogger) SErrorf(format string, args ...any) string {
	return fmt.Sprintf(fmt.Sprintf("%s %s %s\n", ErrorLabel(), DateLabel(), format), args...)
}

func (l defaultLogger) SWarnf(format string, args ...any) string {
	return fmt.Sprintf(fmt.Sprintf("%s %s %s\n", WarnLabel(), DateLabel(), format), args...)
}

func (l defaultLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(l.SInfof(format, args...))
}

func (l defaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(l.SErrorf(format, args...))
}

func (l defaultLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(l.SWarnf(format, args...))
}

func (l defaultLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(l.SErrorf(format, args...))
}

func WithColor(c, msg string) string {
	return fmt.Sprintf("%s%s%s", c, msg, Reset)
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
