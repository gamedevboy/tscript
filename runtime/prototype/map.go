package prototype

import (
    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/runtime_t"
    "tklibs/script/type/function"
    "tklibs/script/type/object"
)

type Map struct {
    *object.Component
    context interface{}
}

func (impl *Map) NativeCall(_, _ interface{}, args ...interface{}) (interface{}, error) {
    capSize := 0
    if len(args) > 0 {
        capSize = int(args[0].(script.Int))
    }
    return impl.context.(runtime.ScriptContext).NewScriptMap(capSize), nil
}

func (impl *Map) GetContext() interface{} {
    return impl.context
}

func (impl *Map) IsScriptFunction() bool {
    return false
}

func (impl *Map) GetScriptRuntimeFunction() runtime_t.Function {
    return nil
}

func (impl *Map) GetNativeRuntimeFunction() runtime_t.NativeFunction {
    return nil
}

func (*Map) Reload() {
}

func (*Map) Invoke(context, this interface{}, _ ...interface{}) interface{} {
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

var _ runtime_t.NativeFunction = &Map{}
var _ script.Function = &Map{}

func NewMapPrototype(ctx interface{}) *Map {
    obj := &Map{
        context: ctx,
    }
    obj.Component = object.NewScriptObject(obj, ctx, 0)

    obj.ScriptSet("forEach", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        f := args[0].(script.Function)

        return this.(script.Map).Foreach(func(key, value interface{}) bool {
            ret := f.Invoke(context, script.Null, key, value)
            if r, ok := ret.(script.Bool); ok && r == false {
                return false
            }
            return true
        })
    }, ctx))

    obj.ScriptSet("length", function.NativeFunctionToValue(func(context interface{}, this interface{}, _ ...interface{}) interface{} {
        return this.(script.Map).Len()
    }, ctx))

    obj.ScriptSet("values", function.NativeFunctionToValue(func(context interface{}, this interface{},
        _ ...interface{}) interface{} {
        thisMap := this.(script.Map)
        array := context.(runtime.ScriptContext).NewScriptArray(int(thisMap.Len()))
        arrayList := array.(script.Array)

        thisMap.Foreach(func(key, value interface{}) bool {
            arrayList.Push(script.ToValue(value))
            return true
        })

        return array
    }, ctx))

    obj.ScriptSet("keys", function.NativeFunctionToValue(func(context interface{}, this interface{},
        _ ...interface{}) interface{} {
        thisMap := this.(script.Map)
        array := context.(runtime.ScriptContext).NewScriptArray(int(thisMap.Len()))
        arrayList := array.(script.Array)

        thisMap.Foreach(func(key, value interface{}) bool {
            arrayList.Push(script.ToValue(key))
            return true
        })

        return array
    }, ctx))

    obj.ScriptSet("containsKey", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        return this.(script.Map).ContainsKey(args[0])
    }, ctx))

    obj.ScriptSet("set", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        if len(args) < 2 {
            return this
        }
        this.(script.Map).Set(args[0], args[1])
        return this
    }, ctx))

    obj.ScriptSet("get", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.Null
        }
        return this.(script.Map).Get(args[0])
    }, ctx))

    obj.ScriptSet("delete", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return this
        }
        this.(script.Map).Delete(args[0])
        return this
    }, ctx))

    return obj
}
