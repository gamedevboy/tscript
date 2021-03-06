package io

import (
    "io/ioutil"
    "os"
    "time"

    "tklibs/script"
    "tklibs/script/runtime/native"
)

type library struct {
    context interface{}
    UnixNow,
    ReadAll, WriteAll native.FunctionType
}

func (*library) GetName() string {
    return "io"
}

func (l *library) SetScriptContext(context interface{}) {
    l.context = context
}

func NewLibrary() *library {
    ret := &library{}
    ret.init()
    return ret
}

func (l *library) init() {
    l.UnixNow = func(_ interface{}, _ interface{}, _ ...interface{}) interface{} {
        return script.Int64(time.Now().UTC().UnixNano() / 1000000)
    }

    l.ReadAll = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.String("")
        }

        if file, err := os.Open(string(args[0].(script.String))); err == nil {
            defer file.Close()
            if buf, err := ioutil.ReadAll(file); err == nil {
                return script.String(buf)
            }
        }

        return script.String("")
    }

    l.WriteAll = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        if len(args) < 2 {
            return script.String("")
        }

        if file, err := os.Create(string(args[0].(script.String))); err == nil {
            defer file.Close()
            if n, err := file.WriteString(string(args[1].(script.String))); err == nil {
                return script.Int(n)
            }
        }

        return script.Int(-1)
    }
}
