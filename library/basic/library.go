package basic

import (
    "fmt"
    "math"
    "strconv"

    "tklibs/script"
    "tklibs/script/runtime/function/native"
    "tklibs/script/runtime/util"
)

type library struct {
    context interface{}
    SetPrototype,
    GetPrototype,
    ToInt,
    ToFloat,
    Print,
    Println,
    TypeOf native.FunctionType
}

func (l *library) SetScriptContext(context interface{}) {
    l.context = context
}

func (*library) GetName() string {
    return ""
}

var Library = &library{}

func init() {
    Library.Print = func(this interface{}, args ...interface{}) interface{} {
        fmt.Print(args...)
        return this
    }

    Library.Println = func(this interface{}, args ...interface{}) interface{} {
        fmt.Println(args...)
        return this
    }

    Library.ToInt = func(this interface{}, args ...interface{}) interface{} {
        switch value := args[0].(type) {
        case script.Int:
            return value
        case script.Float:
            return script.Int(value)
        case script.Float64:
            return script.Int(value)
        case script.Bool:
            if value {
                return script.Int(1)
            }
            return script.Int(0)
        case script.String:
            val, err := strconv.ParseInt(string(value), 10, 64)
            if err != nil {
                return script.Int(0)
            }

            if val > math.MaxInt32 || val < math.MinInt32 {
                return script.Int64(val)
            }

            return script.Int(val)
        default:
            return this
        }
    }

    Library.ToFloat = func(this interface{}, args ...interface{}) interface{} {
        return util.ToScriptFloat(this)
    }

    Library.SetPrototype = func(this interface{}, args ...interface{}) interface{} {
        return script.NullValue
    }

    Library.TypeOf = func(_ interface{}, args ...interface{}) interface{} {
        if len(args) < 1 {
            return script.String("")
        }

        this := args[0]

        if this == script.Null {
            return script.String("null")
        }

        switch this.(type) {
        case script.String:
            return script.String("string")
        case script.Bool:
            return script.String("boolean")
        case script.Int,
            script.Int64,
            script.Float,
            script.Float64:
            return script.String("number")
        case script.Object:
            return script.String("object")
        default:
            return script.String("Unknown type")
        }
    }
}
