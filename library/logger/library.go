package logger

import (
	"tklibs/script"
	"tklibs/script/runtime/native"
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

func (l *library) init() {
	l.Debug = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Debug(args...)
		return this
	}
	l.Debugf = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Debugf(format, args[1:]...)
		} else {
			ScriptLogger().Debug(args...)
		}
		return this
	}

	l.Info = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Info(args...)
		return this
	}
	l.Infof = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Infof(format, args[1:]...)
		} else {
			ScriptLogger().Info(args...)
		}
		return this
	}

	l.Warning = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Warning(args...)
		return this
	}
	l.Warningf = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Warningf(format, args[1:]...)
		} else {
			ScriptLogger().Warning(args...)
		}
		return this
	}

	l.Error = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Error(args...)
		return this
	}
	l.Errorf = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Errorf(format, args[1:]...)
		} else {
			ScriptLogger().Error(args...)
		}
		return this
	}

	l.Fatal = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Fatal(args...)
		return this
	}
	l.Fatalf = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Fatalf(format, args[1:]...)
		} else {
			ScriptLogger().Fatal(args...)
		}
		return this
	}

	l.Panic = func(this interface{}, args ...interface{}) interface{} {
		ScriptLogger().Panic(args...)
		return this
	}
	l.Panicf = func(this interface{}, args ...interface{}) interface{} {
		if format, ok := extract(args); ok {
			ScriptLogger().Panicf(format, args[1:]...)
		} else {
			ScriptLogger().Panic(args...)
		}
		return this
	}
}
