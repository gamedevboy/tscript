package native

import (
    "reflect"

    "tklibs/script/runtime/runtime_t"
)

type ReflectFunction reflect.Value

var _ runtime_t.NativeFunction = ReflectFunction{}

func (impl ReflectFunction) NativeCall(context, this interface{}, args ...interface{}) (interface{}, error) {
    callArgs := make([]reflect.Value, len(args)+2)

    callArgs[0] = reflect.ValueOf(this)

    for idx, value := range args {
        callArgs[idx+1] = reflect.ValueOf(value)
    }

    ret := reflect.Value(impl).Call(callArgs)

    if len(ret) > 0 {
        return ret[0], nil
    }

    return this, nil
}

type FunctionType func(interface{}, interface{}, ...interface{}) interface{}

func (fn FunctionType) NativeCall(context, this interface{}, args ...interface{}) (r interface{}, err error) {
    r = fn(context, this, args...)
    return
}
