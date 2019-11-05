package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(fmt string,v ...interface{})
	Error(v ...interface{})
	Errorf(fmt string,v ...interface{})
	Info(v ...interface{})
	Infof(fmt string,v ...interface{})
	Warning(v ...interface{})
	Warningf(fmt string,v ...interface{})
	Fatal(v ...interface{})
	Fatalf(fmt string,v ...interface{})
	Panic(v ...interface{})
	Panicf(fmt string,v ...interface{})
}

func ScriptLogger() Logger {
	return scriptLogger
}

func SetScriptLogger(l Logger) {
	scriptLogger = l
}

var (
	defaultLogger = &DefaultLogger{Logger: log.New(os.Stdout, "[tscript]", log.LstdFlags),debug:true}
	discardLogger = &DefaultLogger{Logger: log.New(ioutil.Discard, "", 0)}
	scriptLogger = Logger(defaultLogger)
)

const (
	callDepth = 2
)

// DefaultLogger is a default implementation of the Logger interface.
type DefaultLogger struct {
	*log.Logger
	debug bool
}

func (l *DefaultLogger) EnableTimestamps() {
	l.SetFlags(l.Flags() | log.Ldate | log.Ltime)
}

func (l *DefaultLogger) EnableDebug() {
	l.debug = true
}

func prefix(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}

func (l *DefaultLogger) Debug(v ...interface{}) {
	if l.debug {
		_ = l.Output(callDepth, prefix("DEBUG", fmt.Sprint(v...)))
	}
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	if l.debug {
		_ = l.Output(callDepth, prefix("DEBUG", fmt.Sprintf(format, v...)))
	}
}

func (l *DefaultLogger) Info(v ...interface{}) {
	_ = l.Output(callDepth, prefix("INFO", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("INFO", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Error(v ...interface{}) {
	_ = l.Output(callDepth, prefix("ERROR", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("ERROR", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Warning(v ...interface{}) {
	_ = l.Output(callDepth, prefix("WARN", fmt.Sprint(v...)))
}

func (l *DefaultLogger) Warningf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("WARN", fmt.Sprintf(format, v...)))
}

func (l *DefaultLogger) Fatal(v ...interface{}) {
	_ = l.Output(callDepth, prefix("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Fatalf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *DefaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *DefaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}
