package runtime

import "tklibs/script"

type ScriptInterpreter interface {
    InvokeFunction(function, this interface{}, args ...interface{}) interface{}
    InvokeNew(function interface{}, args ...interface{}) interface{}
    GetPC() int
    GetCurrentRegisters() []script.Value
}
