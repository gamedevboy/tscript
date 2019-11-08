package script

import (
    "tklibs/script/runtime/runtime_t"
)

type Function interface {
    Invoke(interface{}, ...interface{}) interface{}
    SetThis(Value)
    GetThis() Value

    GetRuntimeFunction() interface{}

    IsScriptFunction() bool

    GetScriptRuntimeFunction() runtime_t.Function
    GetNativeRuntimeFunction() runtime_t.NativeFunction

    Reload()
    GetRefList() []*Value

    GetFieldByMemberIndex(obj interface{}, index Int) Value
    SetFieldByMemberIndex(obj interface{}, index Int, value Value)
}
