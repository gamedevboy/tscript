package value

import (
    "fmt"
    "math"
    "strconv"
    "strings"

    "tklibs/script"
    "tklibs/script/runtime"
)

func FromInterfaceSlice(array []interface{}) []script.Value {
    ret := make([]script.Value, len(array))

    for i, val := range array {
        switch value := val.(type) {
        case script.Int:
            ret[i].SetInt(value)
        case script.Int64:
            ret[i].SetInterface(val)
        case script.Float:
            ret[i].SetFloat(value)
        case script.Bool:
            ret[i].SetBool(value)
        default:
            ret[i].SetInterface(value)
        }
    }

    return ret
}

func ToInterfaceSlice(array []script.Value) []interface{} {
    ret := make([]interface{}, len(array))

    for i, val := range array {
        switch val.GetType() {
        case script.ValueTypeInt:
            ret[i] = val.GetInt()
        case script.ValueTypeFloat:
            ret[i] = val.GetFloat()
        case script.ValueTypeBool:
            ret[i] = val.GetBool()
        default:
            ret[i] = val.GetInterface()
        }
    }

    return ret
}

func ToJsonString(value interface{}) string {
    switch v := value.(type) {
    case script.Bool:
        return fmt.Sprint(v)
    case script.Int:
        return fmt.Sprint(v)
    case script.Float:
        return fmt.Sprint(v)
    case script.Int64:
        return fmt.Sprint(v)
    case script.Float64:
        return fmt.Sprint(v)
    case script.String:
        return strconv.Quote(string(v))
    case script.Array:
        sb := strings.Builder{}
        sb.WriteRune('[')
        arrayLen := v.Len()
        for i := script.Int(0); i < arrayLen; i++ {
            sb.WriteString(ToJsonString(v.GetElement(i).Get()))
            if i < arrayLen-1 {
                sb.WriteRune(',')
            }
        }
        sb.WriteRune(']')
        return sb.String()
    case script.Map:
        sb := strings.Builder{}
        sb.WriteRune('{')
        l := v.Len()
        i := 0
        v.Foreach(func(key, val interface{}) bool {
            sb.WriteString(fmt.Sprintf("\"%v\":", key.(script.String)))
            sb.WriteString(ToJsonString(val))
            i++
            if i < int(l) {
                sb.WriteRune(',')
            }
            return true
        })

        sb.WriteRune('}')
        return sb.String()
    case runtime.Object:
        sb := strings.Builder{}
        sb.WriteRune('{')
        names := v.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldNames()
        for i, name := range names {
            obj := *v.GetByIndex(i)
            switch v := obj.Get().(type) {
            case script.Function:
            default:
                sb.WriteString(fmt.Sprintf("\"%v\":%v", *name, ToJsonString(v)))
                if i < len(names)-1 {
                    sb.WriteRune(',')
                }
            }
        }
        sb.WriteRune('}')
        return sb.String()
    }

    return "null"
}

func ToScriptString(context, value interface{}) script.String {
    switch dest := value.(type) {
    case script.String:
        return dest
    case script.Int:
        return script.String(strconv.FormatInt(int64(int(dest)), 10))
    case script.Float:
        return script.String(strconv.FormatFloat(float64(dest), 'f', -1, 32))
    case script.Int64:
        return script.String(strconv.FormatInt(int64(dest), 10))
    case script.Float64:
        return script.String(strconv.FormatFloat(float64(dest), 'f', -1, 64))
    case script.Object:
        if dest.GetScriptTypeId() == script.ScriptTypeNull {
            return "null"
        } else {
            funcValue := dest.ScriptGet("toString")
            if _func, ok := funcValue.Get().(script.Function); ok && _func != nil {
                switch retValue := _func.Invoke(context, dest).(type) {
                case script.String:
                    return retValue
                default:
                    panic(fmt.Sprintf("Can not convert '%v'.toString() to String", value))
                }
            }

            return script.String(fmt.Sprintln(dest))
        }
    default:
        return script.String(fmt.Sprint(dest))
    }
}

func ToScriptFloat(value interface{}) script.Float {
    switch dest := value.(type) {
    case script.Float:
        return dest
    case script.Int:
        return script.Float(dest)
    case script.String:
        val, _ := strconv.ParseFloat(string(dest), 32)
        return script.Float(val)
    default:
        return script.Float(math.NaN())
    }
}

func ToScriptInt64(value interface{}) script.Int64 {
    switch dest := value.(type) {
    case script.Int:
        return script.Int64(dest)
    case script.Int64:
        return dest
    case script.String:
        val, _ := strconv.ParseFloat(string(dest), 64)
        return script.Int64(val)
    case script.Bool:
        if dest {
            return 1
        }
    }

    return 0
}
