package prototype

import (
    "fmt"
    "sort"
    "strings"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/runtime_t"
    "tklibs/script/type/function"
    "tklibs/script/type/object"
    "tklibs/script/value"
)

type Array struct {
    *object.Component
    context interface{}
}

func (a *Array) NativeCall(_ interface{}, args ...interface{}) interface{} {
    capSize := 0
    if len(args) > 0 {
        capSize = int(args[0].(script.Int))
    }
    return a.context.(runtime.ScriptContext).NewScriptArray(capSize)
}

func (a *Array) GetContext() interface{} {
    return a.context
}

func (a *Array) IsScriptFunction() bool {
    return false
}

func (a *Array) GetScriptRuntimeFunction() runtime_t.Function {
    return nil
}

func (a *Array) GetNativeRuntimeFunction() runtime_t.NativeFunction {
    return nil
}

func (*Array) Reload() {

}

func (a *Array) Invoke(this interface{}, args ...interface{}) interface{} {
    return this
}

func (a *Array) SetThis(script.Value) {
}

func (a *Array) GetThis() script.Value {
    return script.Value{}
}

func (a *Array) GetRuntimeFunction() interface{} {
    return a
}

func (a *Array) GetRefList() []*script.Value {
    return nil
}

func (a *Array) GetFieldByMemberIndex(obj interface{}, index script.Int) script.Value {
    return script.Value{}
}

func (a *Array) SetFieldByMemberIndex(obj interface{}, index script.Int, value script.Value) {
}

var _ script.Function = &Array{}
var _ runtime_t.NativeFunction = &Array{}

func NewArrayPrototype(ctx interface{}) *Array {
    obj := &Array{
        context: ctx,
    }
    obj.Component = object.NewScriptObject(obj, ctx, 0)

    obj.ScriptSet("push", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        return this.(script.Array).Push(value.FromInterfaceSlice(args)...)
    }, ctx))

    obj.ScriptSet("pop", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return this.(script.Array).Pop()
    },ctx))

    obj.ScriptSet("unshift", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        return this.(script.Array).Unshift(value.FromInterfaceSlice(args)...)
    },ctx))

    obj.ScriptSet("shift", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return this.(script.Array).Shift()
    },ctx))

    obj.ScriptSet("length", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(this.(script.Array).Len())
    },ctx))

    obj.ScriptSet("clear", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        this.(script.Array).Clear()
        return this
    },ctx))

    obj.ScriptSet("find", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        f := args[0].(script.Function)

        for i := script.Int(0); i < array.Len(); i++ {
            element := array.GetElement(i).Get()
            ret := f.Invoke(nil, element)
            if ret != nil {
                if r, ok := ret.(script.Bool); ok {
                    if r {
                        return element
                    }
                }
            }
        }

        return script.Null
    },ctx))

    obj.ScriptSet("sort", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        f := args[0].(script.Function)

        sort.Slice(array.GetSlice(), func(i,j int) bool {
            ret := f.Invoke(nil, array.GetElement(script.Int(i)).Get(), array.GetElement(script.Int(j)).Get())
            return ret == script.Bool(true)
        })

        return script.Null
    },ctx))

    obj.ScriptSet("findIndex", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        f := args[0].(script.Function)

        for i := script.Int(0); i < array.Len(); i++ {
            element := array.GetElement(i).Get()
            ret := f.Invoke(nil, element)
            if ret != nil {
                if r, ok := ret.(script.Bool); script.Bool(ok) && r {
                    return i
                }
            }
        }

        return script.Int(-1)
    },ctx))

    obj.ScriptSet("map", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        ret := ctx.(runtime.ScriptContext).NewScriptArray(int(array.Len()))
        f := args[0].(script.Function)
        retArray := ret.(script.Array)
        v := script.Value{}
        for i := script.Int(0); i < array.Len(); i++ {
            retArray.SetElement(i, v.Set(f.Invoke(nil, array.GetElement(i).Get())))
        }
        return ret
    },ctx))

    obj.ScriptSet("forEach", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        f := args[0].(script.Function)
        i := script.Int(0)
        for ; i < array.Len(); i++ {
            ret := f.Invoke(nil, array.GetElement(i).Get(), i)
            if ret != nil {
                if r, ok := ret.(script.Bool); ok && r == false {
                    return i
                }
            }
        }
        return i
    },ctx))

    obj.ScriptSet("removeAt", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        return array.RemoveAt(args[0].(script.Int))
    },ctx))

    obj.ScriptSet("join", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.String("")
        }

        array := this.(script.Array)
        a := make([]string, array.Len())
        for i := script.Int(0); i < array.Len(); i++ {
            a[i] = fmt.Sprint(array.GetElement(i).Get())
        }

        return script.String(strings.Join(a, string(args[0].(script.String))))
    },ctx))

    return obj
}
