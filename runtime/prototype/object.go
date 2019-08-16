package prototype

import (
    "fmt"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/type/object"
)

type Object struct {
    script.ComponentType
    prototype interface{}
}

func (impl *Object) GetObjectPrototype() interface{} {
    return impl.prototype
}

func (impl *Object) Init() {
    obj := impl.prototype.(script.Object)

    obj.ScriptSet("toString", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        switch str := this.(type) {
        case fmt.Stringer:
            return script.String(str.String())
        default:
            return script.String(fmt.Sprint(this))
        }
    }).ToValue(impl.GetOwner()))
}

func NewObjectPrototype(ctx interface{}) *Object {
    ret := &Object{
        ComponentType: script.MakeComponentType(ctx),
    }

    n := &struct {
        *object.Component
    }{}
    n.Component = object.NewScriptObjectWithRuntimePrototype(n, ctx.(runtime.ScriptContext).GetRootRuntimeType(), nil, 2)
    ret.prototype = n
    return ret
}
