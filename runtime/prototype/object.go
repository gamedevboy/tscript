package prototype

import (
    "fmt"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/type/function"
    "tklibs/script/type/object"
)

type Object struct {
    script.ComponentType
    prototype interface{}
}

func (impl *Object) GetObjectPrototype() interface{} {
    return impl.prototype
}

func (impl *Object) InitPrototype() {
    obj := impl.prototype.(script.Object)

    obj.ScriptSet("toString", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        switch str := this.(type) {
        case fmt.Stringer:
            return script.String(str.String())
        default:
            return script.String(fmt.Sprint(this))
        }
    },impl.GetOwner()))

    obj.ScriptSet("setPrototype", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        this.(runtime.Object).SetPrototype(script.InterfaceToValue(args[0]))
        return this
    },impl.GetOwner()))

    obj.ScriptSet("getPrototype", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return this.(runtime.Object).GetPrototype()
    },impl.GetOwner()))
}

func NewObjectPrototype(ctx interface{}) *Object {
    ret := &Object{
        ComponentType: script.MakeComponentType(ctx),
    }

    n := &struct {
        *object.Component
    }{}
    n.Component = object.NewScriptObjectWithRuntimePrototype(n, ctx.(runtime.ScriptContext).GetRootRuntimeType(), ctx,nil, 2)
    ret.prototype = n
    return ret
}
