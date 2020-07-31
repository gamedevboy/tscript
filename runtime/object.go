package runtime

import "tklibs/script"

type Object interface {
    GetRuntimeTypeInfo() interface{}
    GetByIndex(int) *script.Value
    SetByIndex(int, script.Value) *script.Value
    GetPrototype() script.Value
    SetPrototype(value script.Value)
}
