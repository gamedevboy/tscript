package prototype

import (
    "fmt"
    "strings"
    "unicode/utf8"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/type/function"
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

    obj.ScriptSet("format", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        return script.String(fmt.Sprintf(strings.Replace(string(args[0].(script.String)), "{}", "%v", -1), args[1:]...))
    }, ctx))

    obj.ScriptSet("length", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(utf8.RuneCountInString(string(this.(script.String))))
    }, ctx))

    obj.ScriptSet("replace", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        return script.String(strings.Replace(string(this.(script.String)), string(args[0].(script.String)), string(args[1].(script.String)),
            -1))
    }, ctx))

    obj.ScriptSet("length", function.NativeFunctionToValue(func(this interface{}, _ ...interface{}) interface{} {
        return script.Int(utf8.RuneCountInString(string(this.(script.String))))
    }, ctx))

    obj.ScriptSet("indexOf", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        index := strings.Index(string(this.(script.String)), string(this.(script.String)))
        if index == - 1 {
            return script.Int(-1)
        }
        start := []byte(this.(script.String))[:index]
        return script.Int(utf8.RuneCount(start))
    }, ctx))

    obj.ScriptSet("lastIndexOf", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        index := strings.LastIndex(string(this.(script.String)), string(this.(script.String)))
        if index == - 1 {
            return script.Int(-1)
        }
        start := []byte(this.(script.String))[:index]
        return script.Int(utf8.RuneCount(start))
    }, ctx))

    obj.ScriptSet("join", function.NativeFunctionToValue(func(_ interface{}, args ...interface{}) interface{} {
        if len(args) < 2 {
            return script.String("")
        }

        array := args[1].(script.Array)
        a := make([]string, array.Len())
        for i := script.Int(0); i < array.Len(); i++ {
            a[i] = fmt.Sprint(array.GetElement(i).Get())
        }

        return script.String(strings.Join(a, string(args[0].(script.String))))
    }, ctx))

    obj.ScriptSet("sub", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
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
    }, ctx))

    obj.ScriptSet("+", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
        argLen := len(args)
        if argLen < 1 {
            return this
        }

        return script.String(append([]rune(this.(script.String)), []rune(args[0].(script.String))...))
    }, ctx))

    return ret
}
