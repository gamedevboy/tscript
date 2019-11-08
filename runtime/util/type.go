package util

import (
    "fmt"

    "tklibs/script"
)

func ToScriptString(value interface{}) script.String {
    switch dest := value.(type) {
    case script.String:
        return dest
    case script.Object:
        funcValue := dest.ScriptGet("toString")
        switch retValue := funcValue.GetFunction().Invoke(dest).(type) {
        case script.String:
            return retValue
        default:
            panic(fmt.Sprintf("Can not convert '%v'.toString() to String", value))
        }
    default:
        return script.String(fmt.Sprint(dest))
    }
}

func ToScriptInt(value interface{}) script.Int {
    switch dest := value.(type) {
    case *script.Value:
        return dest.GetInt()
    default:
        panic(fmt.Sprintf("Can not convert '%v' to Int", value))
    }
}

func ToScriptBool(value interface{}) script.Bool {
    if value == script.Null {
        return false
    }

    switch dest := value.(type) {
    case script.Bool:
        return dest
    case script.Float:
        return dest != 0
    case script.Int:
        return dest != 0
    case script.Object:
        return true
    default:
        panic("")
    }
}

func ToScriptFloat(value interface{}) script.Float {
    switch dest := value.(type) {
    case script.Float:
        return dest
    case script.Int:
        return script.Float(dest)
    default:
        panic("can")
    }
}
