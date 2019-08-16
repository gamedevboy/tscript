package io

import (
    "io/ioutil"
    "os"
    "time"

    "tklibs/script"
    "tklibs/script/runtime/function/native"
)

type library struct {
    context interface{}
    UnixNow,
    ReadAll native.FunctionType
}

func (*library) GetName() string {
    return "io"
}

func (l *library) SetScriptContext(context interface{}) {
    l.context = context
}

var Library = &library{}

func init() {
    Library.UnixNow = func(_ interface{}, _ ...interface{}) interface{} {
        return script.Int64(time.Now().UTC().UnixNano() / 1000000)
    }

    Library.ReadAll = func(this interface{}, args ...interface{}) interface{} {
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
}
