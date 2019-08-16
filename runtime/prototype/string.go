package prototype

import (
    "fmt"
    "strings"
    "unicode/utf8"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
)

type String struct {
    prototype interface{}
}

func (impl *String) GetStringPrototype() interface{} {
    return impl.prototype
}

func NewStringPrototype(ctx interface{}) *String {
    ret := &String{}
    ret.prototype = ctx.(runtime.ScriptContext).NewScriptObject(0)

    obj := ret.prototype.(script.Object)

    obj.ScriptSet("format", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        return script.String(fmt.Sprintf(strings.Replace(string(args[0].(script.String)), "{}", "%v", -1), args[1:]...))
    }).ToValue(ctx))

    obj.ScriptSet("length", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(utf8.RuneCountInString(string(this.(script.String))))
    }).ToValue(ctx))

    obj.ScriptSet("replace", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        return script.String(strings.Replace(string(this.(script.String)), string(args[0].(script.String)), string(args[1].(script.String)),
            -1))
    }).ToValue(ctx))

    obj.ScriptSet("length", native.FunctionType(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(utf8.RuneCountInString(string(this.(script.String))))
    }).ToValue(ctx))

    obj.ScriptSet("indexOf", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        index := strings.Index(string(this.(script.String)), string(this.(script.String)))
        if index == - 1 {
            return script.Int(-1)
        }
        start := []byte(this.(script.String))[:index]
        return script.Int(utf8.RuneCount(start))
    }).ToValue(ctx))

    obj.ScriptSet("lastIndexOf", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        index := strings.LastIndex(string(this.(script.String)), string(this.(script.String)))
        if index == - 1 {
            return script.Int(-1)
        }
        start := []byte(this.(script.String))[:index]
        return script.Int(utf8.RuneCount(start))
    }).ToValue(ctx))

    obj.ScriptSet("join", native.FunctionType(func(_ interface{}, args ...interface{}) interface{} {
        if len(args) < 2 {
            return script.String("")
        }

        array := args[1].(script.Array)
        a := make([]string, array.Len())
        for i := script.Int(0); i < array.Len(); i++ {
            a[i] = fmt.Sprint(array.GetElement(i).Get())
        }

        return script.String(strings.Join(a, string(args[0].(script.String))))
    }).ToValue(ctx))

    obj.ScriptSet("sub", native.FunctionType(func(this interface{}, args ...interface{}) interface{} {
        argLen := len(args)
        if argLen < 1 {
            return this
        }

        runes := []rune(string(this.(script.String)))

        startIndex := args[0].(script.Int)
        if argLen < 2 {
            return script.String(runes[startIndex:])
        }

        count := args[1].(script.Int)
        return script.String(runes[startIndex : startIndex+count])
    }).ToValue(ctx))

    return ret
}
