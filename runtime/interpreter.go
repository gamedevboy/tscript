package runtime

import "tklibs/script"

type ScriptInterpreter interface {
	InvokeFunction(function, context, this interface{}, args ...interface{}) interface{}
	InvokeNew(function, context interface{}, args ...interface{}) interface{}
	GetPC() int
	GetCurrentRegisters() []script.Value
}
