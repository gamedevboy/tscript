package native

import (
    "reflect"

    "tklibs/script/runtime/runtime_t"
)

type ReflectFunction reflect.Value

var _ runtime_t.NativeFunction = ReflectFunction{}

func (impl ReflectFunction) NativeCall(this interface{}, args ...interface{}) interface{} {
    callArgs := make([]reflect.Value, len(args)+1)

    callArgs[0] = reflect.ValueOf(this)

    for idx, value := range args {
        callArgs[idx+1] = reflect.ValueOf(value)
    }

    ret := reflect.Value(impl).Call(callArgs)

    if len(ret) > 0 {
        return ret[0]
    }

    return this
}

type FunctionType func(interface{}, ...interface{}) interface{}

func (fn FunctionType) NativeCall(this interface{}, args ...interface{}) interface{} {
    return fn(this, args...)
}