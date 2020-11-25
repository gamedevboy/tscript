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

    obj.ScriptSet("toString", function.NativeFunctionToValue(func(context interface{}, this interface{}, _ ...interface{}) interface{} {
        switch str := this.(type) {
        case fmt.Stringer:
            return script.String(str.String())
        default:
            return script.String(fmt.Sprint(this))
        }
    },impl.GetOwner()))

    obj.ScriptSet("setPrototype", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        this.(runtime.Object).SetPrototype(script.InterfaceToValue(args[0]))
        return this
    },impl.GetOwner()))

    obj.ScriptSet("getPrototype", function.NativeFunctionToValue(func(context interface{}, this interface{}, _ ...interface{}) interface{} {
        return this.(runtime.Object).GetPrototype()
    },impl.GetOwner()))

    obj.ScriptSet("forEach", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        scriptObj := this.(runtime.Object)
        rt := scriptObj.GetRuntimeTypeInfo().(runtime.TypeInfo)
        f := args[0].(script.Function)

        for i, name := range  rt.GetFieldNames() {
            ret := f.Invoke(context, script.Null, script.String(name), scriptObj.GetByIndex(i).Get(), i)
            if ret != nil {
                if r, ok := ret.(script.Bool); ok && r == false {
                    return i
                }
            }
        }

        return this
    },impl.GetOwner()))
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
