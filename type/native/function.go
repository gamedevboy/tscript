package native

import (
    "tklibs/script"
    "tklibs/script/runtime/native"
)

type function struct {
    funcValue native.ReflectFunction
    name      string
}

func (f *function) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeFunction
}

func (f *function) ScriptSet(string, *script.Value) {
    panic("not supported")
}

func (f *function) ScriptGet(string) *script.Value {
    panic("not supported")
}

func (f *function) Invoke(this interface{}, args ...interface{}) interface{} {
    return f.funcValue.NativeCall(this, args...)
}

func (f *function) GetRuntimeFunction() interface{} {
    return f.funcValue
}

func (*function) GetRefList() []*script.Value {
    panic("not supported")
}

func (*function) GetFieldByMemberIndex(obj interface{}, index script.Int) *script.Value {
    panic("not supported")
}

func (*function) SetFieldByMemberIndex(obj interface{}, index script.Int, value *script.Value) {
    panic("not supported")
}

var _ script.Function = &function{}
var _ script.Object = &function{}
