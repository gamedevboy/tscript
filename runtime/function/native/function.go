package native

import (
    "reflect"

    "tklibs/script"
    "tklibs/script/type/function"
)

type ReflectFunction reflect.Value

var _ Function = ReflectFunction{}

func (impl ReflectFunction) NativeCall(this interface{}, args ...interface{}) script.Value {
    callArgs := make([]reflect.Value, len(args)+1)

    callArgs[0] = reflect.ValueOf(this)

    for idx, value := range args {
        callArgs[idx+1] = reflect.ValueOf(value)
    }

    ret := reflect.Value(impl).Call(callArgs)
    r := script.Value{}
    if len(ret) > 0 {
        return r.Set(ret[0])
    }

    return r.Set(this)
}

type FunctionType func(interface{}, ...interface{}) interface{}

func (fn FunctionType) NativeCall(this interface{}, args ...interface{}) script.Value {
    r := script.Value{}
    return r.Set(fn(this, args...))
}

func (fn FunctionType) ToValue(context interface{}) script.Value {
    return script.InterfaceToValue(NewNativeFunction(fn, context))
}

type Function interface {
    NativeCall(interface{}, ...interface{}) script.Value
}

func NewNativeFunction(f, ctx interface{}) interface{} {
    nativeFunction := &struct {
        *function.Component
    }{}

    switch _func := f.(type) {
    case FunctionType:
        nativeFunction.Component = function.NewScriptFunction(nativeFunction, _func, ctx)
    case reflect.Value:
        nativeFunction.Component = function.NewScriptFunction(nativeFunction, ReflectFunction(_func), ctx)
    default:
        panic("")
    }

    return nativeFunction
}
