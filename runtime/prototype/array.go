package prototype

import (
    "fmt"
    "strings"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/type/object"
    "tklibs/script/value"
)

type Array struct {
    *object.Component
    context interface{}
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

func (a *Array) New(args ...interface{}) interface{} {
    capSize := 0
    if len(args) > 0 {
        capSize = int(args[0].(script.Int))
    }
    return a.context.(runtime.ScriptContext).NewScriptArray(capSize)
}

var _ script.Function = &Array{}
var _ native.Type = &Array{}

func NewArrayPrototype(ctx interface{}) *Array {
    obj := &Array{
        context: ctx,
    }
    obj.Component = object.NewScriptObject(obj, ctx, 0)

    obj.ScriptSet("push", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        return this.(script.Array).Push(value.FromInterfaceSlice(args)...)
    }).ToValue(ctx))

    obj.ScriptSet("pop", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return this.(script.Array).Pop()
    }).ToValue(ctx))

    obj.ScriptSet("unshift", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        return this.(script.Array).Unshift(value.FromInterfaceSlice(args)...)
    }).ToValue(ctx))

    obj.ScriptSet("shift", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return this.(script.Array).Shift()
    }).ToValue(ctx))

    obj.ScriptSet("length", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(this.(script.Array).Len())
    }).ToValue(ctx))

    obj.ScriptSet("find", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
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
    }).ToValue(ctx))

    obj.ScriptSet("findIndex", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
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
    }).ToValue(ctx))

    obj.ScriptSet("map", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        ret := ctx.(runtime.ScriptContext).NewScriptArray(int(array.Len()))
        f := args[0].(script.Function)
        retArray := ret.(script.Array)
        v := script.Value{}
        for i := script.Int(0); i < array.Len(); i++ {
            retArray.SetElement(i, v.Set(f.Invoke(nil, array.GetElement(i).Get())))
        }
        return ret
    }).ToValue(ctx))

    obj.ScriptSet("forEach", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
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
    }).ToValue(ctx))

    obj.ScriptSet("removeAt", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        array := this.(script.Array)
        return array.RemoveAt(args[0].(script.Int))
    }).ToValue(ctx))

    obj.ScriptSet("join", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {            return script.String("")
        }

        array := this.(script.Array)
        a := make([]string, array.Len())
        for i := script.Int(0); i < array.Len(); i++ {
            a[i] = fmt.Sprint(array.GetElement(i).Get())
        }

        return script.String(strings.Join(a, string(args[0].(script.String))))
    }).ToValue(ctx))

    return obj
}
