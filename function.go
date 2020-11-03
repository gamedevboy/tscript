package script

import (
    "tklibs/script/runtime/runtime_t"
)

type Function interface {
    Invoke(context, this interface{}, args ...interface{}) interface{}
    SetThis(Value)
    GetThis() Value

    IsScriptFunction() bool

    GetScriptRuntimeFunction() runtime_t.Function
    GetNativeRuntimeFunction() runtime_t.NativeFunction

    GetRefList() []*Value

    GetFieldByMemberIndex(obj interface{}, index Int) Value
    SetFieldByMemberIndex(obj interface{}, index Int, value Value)
    GetContext() interface{}
}
