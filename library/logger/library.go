package logger

import (
	"tklibs/script"
	"tklibs/script/library/debug"
	"tklibs/script/runtime"
	"tklibs/script/runtime/function/native"
)

type library struct {
	context interface{}
	Debug,
	Debugf,
	Info,
	Infof,
	Warning,
	Warningf,
	Error,
	Errorf,
	Fatal,
	Fatalf,
	Panic,
	Panicf native.FunctionType
}

func (*library) GetName() string {
	return "logger"
}

func (l *library) SetScriptContext(context interface{}) {
	l.context = context
}

func NewLibrary() *library {
	ret := &library{}
	ret.init()
	return ret
}
func extract(args []interface{}) (format string, ok1 bool) {
	if len(args) < 2 {
		return
	}
	fi := args[0]
	if ss, ok := fi.(script.String); ok {
		ok1 = true
		format = string(ss)
		return
	}
	return
}

func logFunc(withFmt bool, log func(v ...interface{}), logf func(fmt string, v ...interface{})) func(this interface{}, args ...interface{}) interface{} {
	if withFmt {
		return func(this interface{}, args ...interface{}) interface{} {
			if format, ok := extract(args); ok {
				logf(format, args[1:]...)
			} else {
				log(args...)
			}
			return this
		}
	}
	return func(this interface{}, args ...interface{}) interface{} {
		log(args...)
		return this
	}
}

func (l *library) init() {
	// todo the logger implementer need 'runtime.ScriptContext' to get CallInfo
	_ = func(context interface{}) *debug.CallInfo {
		return debug.GetCallInfo(context.(runtime.ScriptContext))
	}

	l.Debug = logFunc(false, ScriptLogger().Debug, ScriptLogger().Debugf)
	l.Debugf = logFunc(true, ScriptLogger().Debug, ScriptLogger().Debugf)

	l.Info = logFunc(false, ScriptLogger().Info, ScriptLogger().Infof)
	l.Infof = logFunc(true, ScriptLogger().Info, ScriptLogger().Infof)

	l.Warning = logFunc(false, ScriptLogger().Warning, ScriptLogger().Warningf)
	l.Warningf = logFunc(true, ScriptLogger().Warning, ScriptLogger().Warningf)

	l.Error = logFunc(false, ScriptLogger().Error, ScriptLogger().Errorf)
	l.Errorf = logFunc(true, ScriptLogger().Error, ScriptLogger().Errorf)

	l.Fatal = logFunc(false, ScriptLogger().Fatal, ScriptLogger().Fatalf)
	l.Fatalf = logFunc(true, ScriptLogger().Fatal, ScriptLogger().Fatalf)

	l.Panic = logFunc(false, ScriptLogger().Panic, ScriptLogger().Panicf)
	l.Panicf = logFunc(true, ScriptLogger().Panic, ScriptLogger().Panicf)
}
