package basic

import (
    "fmt"
    "math"
    "strconv"

    "tklibs/script"
    "tklibs/script/runtime/native"
    "tklibs/script/value"
)

type library struct {
    context interface{}
    SetPrototype,
    GetPrototype,
    ToInt,
    ToInt64,
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

func NewLibrary() *library {
    ret := &library{}
    ret.init()
    return ret
}

func (l *library) init() {
    l.Print = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        fmt.Print(args...)
        return this
    }

    l.Println = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        fmt.Println(args...)
        return this
    }

    l.ToInt64 = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        switch value := args[0].(type) {
        case script.Int:
            return script.Int64(value)
        case script.Int64:
            return value
        case script.Float:
            return script.Int64(value)
        case script.Float64:
            return script.Int64(value)
        case script.Bool:
            if value {
                return script.Int64(1)
            }
            return script.Int64(0)
        case script.String:
            val, err := strconv.ParseInt(string(value), 10, 64)
            if err != nil {
                return script.Int64(0)
            }
            return script.Int64(val)
        default:
            return this
        }
    }

    l.ToInt = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        switch value := args[0].(type) {
        case script.Int:
            return value
        case script.Int64:
            return script.Int(value)
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

    l.ToFloat = func(_ interface{}, _ interface{}, args ...interface{}) interface{} {
        return value.ToScriptFloat(args[0])
    }

    l.SetPrototype = func(context interface{}, this interface{}, args ...interface{}) interface{} {
        return script.NullValue
    }

    l.TypeOf = func(context interface{}, _ interface{}, args ...interface{}) interface{} {
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
