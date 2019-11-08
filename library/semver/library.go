package semver

import (
	"tklibs/script"
	"tklibs/script/runtime/function/native"
)

type library struct {
	context interface{}
	IsValid,
	Canonical,
	MajorMinor,
	Prerelease,
	Build,
	Compare,
	Max,
	Major native.FunctionType
}

func (*library) GetName() string {
	return "semver"
}

func (l *library) SetScriptContext(context interface{}) {
	l.context = context
}

func NewLibrary() *library {
	ret := &library{}
	ret.init()
	return ret
}

func (l *library) init() {
	l.IsValid = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return false
		}
		return IsValid(string(args[0].(script.String)))
	}

	l.Canonical = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}
		return Canonical(string(args[0].(script.String)))
	}

	l.MajorMinor = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}
		return MajorMinor(string(args[0].(script.String)))
	}

	l.Prerelease = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}
		return Prerelease(string(args[0].(script.String)))
	}

	l.Build = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}
		return Build(string(args[0].(script.String)))
	}

	l.Compare = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 2 {
			return 0
		}
		return Compare(string(args[0].(script.String)), string(args[1].(script.String)))
	}

	l.Max = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return ""
		}
		if len(args) < 2 {
			return string(args[0].(script.String))
		}
		return script.String(Max(string(args[0].(script.String)), string(args[1].(script.String))))
	}

	l.Major = func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return 0
		}
		return Major(string(args[0].(script.String)))
	}
}
