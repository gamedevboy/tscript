package value

import (
    "fmt"
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

func ToJsonString(value script.Value) string {
    switch value.GetType() {
    case script.ValueTypeBool:
        v := value.GetBool()
        if v {
            return "true"
        } else {
            return "false"
        }
    case script.ValueTypeInt:
        return fmt.Sprint(value.GetInt())
    case script.ValueTypeFloat:
        return fmt.Sprint(value.GetFloat())
    case script.ValueTypeInterface:
        switch v := value.GetInterface().(type) {
        case script.Int64:
            return fmt.Sprint(v)
        case script.Float64:
            return fmt.Sprint(v)
        case script.String:
            return fmt.Sprintf("\"%v\"", v)
        case script.Array:
            sb := strings.Builder{}
            sb.WriteRune('[')
            arrayLen := v.Len()
            for i := script.Int(0); i < arrayLen; i++ {
                sb.WriteString(ToJsonString(v.GetElement(i)))
                if i < arrayLen-1 {
                    sb.WriteRune(',')
                }
            }
            sb.WriteRune(']')
            return sb.String()
        case runtime.Object:
            sb := strings.Builder{}
            sb.WriteRune('{')
            names := v.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldNames()
            for i, name := range names {
                obj := v.GetByIndex(i)
                switch obj.Get().(type) {
                case script.Function:
                default:
                    sb.WriteString(fmt.Sprintf("\"%v\":%v", name, ToJsonString(obj)))
                    if i < len(names)-1 {
                        sb.WriteRune(',')
                    }
                }
            }
            sb.WriteRune('}')
            return sb.String()
        }
    }

    return "null"
}