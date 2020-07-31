package util

import (
	"fmt"
	"math"
	"strconv"

	"tklibs/script"
)

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
	case script.String:
		val, _ := strconv.ParseFloat(string(dest), 32)
		return script.Float(val)
	default:
		return script.Float(math.NaN())
	}
}
