package math

import (
    "math/rand"
    "strconv"

    "tklibs/script"
    "tklibs/script/runtime/function/native"
)

type library struct {
    context interface{}
    ToInt,
    Rand native.FunctionType
}

func (*library) GetName() string {
    return "math"
}

func (l *library) SetScriptContext(context interface{}) {
    l.context = context
}

var Library = &library{}

func init() {
    Library.ToInt = func(this interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.Null
        }

        switch v := args[0].(type) {
        case script.Int:
            return v
        case script.Float:
            return script.Int(v)
        case script.Bool:
            if v {
                return 1
            }
            return 0
        case script.String:
            val, _ := strconv.Atoi(string(v))
            return script.Int(val)
        default:
            return 0
        }
    }

    Library.Rand = func(this interface{}, args ...interface{}) interface{} {
        return script.Float(rand.Float32())
    }
}
