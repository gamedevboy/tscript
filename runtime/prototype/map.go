package prototype

import (
    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
)

type Map struct {
    prototype interface{}
}

func (impl *Map) GetMapPrototype() interface{} {
    return impl.prototype
}

func NewMapPrototype(ctx interface{}) *Map {
    ret := &Map{}
    ret.prototype = ctx.(runtime.ScriptContext).NewScriptObject(0)

    obj := ret.prototype.(script.Object)

    obj.ScriptSet("forEach", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        f := args[0].(script.Function)

        return this.(script.Map).Foreach(func(key, value interface{}) bool {
            ret := f.Invoke(nil, key, value)
            if r, ok := ret.(script.Bool); ok && r == false {
                return false
            }
            return true
        })
    }).ToValue(ctx))

    obj.ScriptSet("length", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return this.(script.Map).Len()
    }).ToValue(ctx))

    obj.ScriptSet("containsKey", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        v := script.Value{}
        return this.(script.Map).ContainsKey(v.Set(args[0]))
    }).ToValue(ctx))

    obj.ScriptSet("delete", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return this
        }
        v := script.Value{}
        this.(script.Map).Delete(v.Set(args[0]))
        return this
    }).ToValue(ctx))

    return ret
}
