package prototype

import (
    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/type/object"
)

type Map struct {
    *object.Component
    context interface{}
}

func (*Map) Reload() {
}

func (*Map) Invoke(this interface{}, _ ...interface{}) interface{} {
    return this
}

func (impl *Map) SetThis(script.Value) {
}

func (impl *Map) GetThis() script.Value {
    return script.Value{}
}

func (impl *Map) GetRuntimeFunction() interface{} {
    return impl
}

func (impl *Map) GetRefList() []*script.Value {
    return nil
}

func (impl *Map) GetFieldByMemberIndex(obj interface{}, index script.Int) script.Value {
    return script.Value{}
}

func (impl *Map) SetFieldByMemberIndex(obj interface{}, index script.Int, value script.Value) {
}

var _ native.Type = &Map{}

func (impl *Map) New(args ...interface{}) interface{} {
    capSize := 0
    if len(args) > 0 {
        capSize = int(args[0].(script.Int))
    }
    return impl.context.(runtime.ScriptContext).NewScriptMap(capSize)
}

func NewMapPrototype(ctx interface{}) *Map {
    obj := &Map{
        context: ctx,
    }
    obj.Component = object.NewScriptObject(obj, ctx, 0)

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
        return this.(script.Map).ContainsKey(args[0])
    }).ToValue(ctx))

    obj.ScriptSet("set", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 2 {
            return this
        }
        this.(script.Map).Set(args[0], args[1])

        return this
    }).ToValue(ctx))

    obj.ScriptSet("get", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.Null
        }
        return this.(script.Map).Get(args[0])
    }).ToValue(ctx))

    obj.ScriptSet("delete", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return this
        }
        this.(script.Map).Delete(args[0])
        return this
    }).ToValue(ctx))

    return obj
}
