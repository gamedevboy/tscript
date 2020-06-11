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
	_defaultLogger = &defaultLogger{Logger: log.New(os.Stdout, "[script]", log.LstdFlags)}
	discardLogger  = &defaultLogger{Logger: log.New(ioutil.Discard, "", 0)}
	scriptLogger   = Logger(_defaultLogger)
)

const (
	callDepth = 2
)

// defaultLogger is a default implementation of the Logger interface.
type defaultLogger struct { *log.Logger}

func prefix(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}

func (l *defaultLogger) Debug(v ...interface{}) {
	_ = l.Output(callDepth, prefix("DEBUG", fmt.Sprint(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	_ = l.Output(callDepth, prefix("INFO", fmt.Sprint(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("INFO", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	_ = l.Output(callDepth, prefix("ERROR", fmt.Sprint(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("ERROR", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Warning(v ...interface{}) {
	_ = l.Output(callDepth, prefix("WARN", fmt.Sprint(v...)))
}

func (l *defaultLogger) Warningf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("WARN", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	_ = l.Output(callDepth, prefix("FATAL", fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	_ = l.Output(callDepth, prefix("FATAL", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}
