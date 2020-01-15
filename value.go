package script

import (
    "fmt"
    "math"
    "strconv"
    "unsafe"
)

const (
    ValueTypeInterface uint16 = iota // must be zero here
    ValueTypeBool
    ValueTypeInt
    ValueTypeFloat
)

type pointer = unsafe.Pointer

type Value struct {
    pointer *Interface
}

func (value *Value) GetType() uint16 {
    return uint16(*(*uint64)(pointer(value)) >> 48)
}

func (value Value) Get() interface{} {
    switch value.GetType() {
    case ValueTypeInt:
        return value.GetInt()
    case ValueTypeFloat:
        return value.GetFloat()
    case ValueTypeBool:
        return value.GetBool()
    default:
        return value.GetInterface()
    }
}

func (value *Value) Set(v interface{}) Value {
    switch val := v.(type) {
    case bool:
        value.SetBool(Bool(val))
    case int:
        value.SetInt(Int(val))
    case int64:
        value.SetInt64(Int64(val))
    case int32:
        value.SetInt(Int(val))
    case uint32:
        value.SetInt64(Int64(val))
    case int8:
        value.SetInt(Int(val))
    case uint8:
        value.SetInt(Int(val))
    case float32:
        value.SetFloat(Float(val))
    case float64:
        value.SetFloat64(Float64(val))
    case string:
        value.SetInterface(String(val))
    case Int:
        value.SetInt(val)
    case Float:
        value.SetFloat(val)
    case Bool:
        value.SetBool(val)
    default:
        value.SetInterface(val)
    }

    return *value
}

func (value Value) Interface() *Interface {
    return value.pointer
}

func (value Value) GetInterface() interface{} {
    if value.pointer == nil {
        return nil
    }
    return value.pointer.Get()
}

func InterfaceToValue(val interface{}) Value {
    value := Value{}
    value.SetInterface(val)
    return value
}

func (value *Value) SetInterface(v interface{}) {
    if v != nil {
        value.pointer = &Interface{}
        value.pointer.Set(v)
    } else {
        value.pointer = nil
    }
}

func (value *Value) GetInt() Int {
    return *(*Int)(pointer(value))
}

func (value *Value) SetInt(v Int) {
    *(*Int)(pointer(value)) = v
    *(*uint32)(pointer(uintptr(pointer(value)) + 4)) = uint32(ValueTypeInt) << 16
}

func (value *Value) GetFloat() Float {
    return *(*Float)(pointer(value))
}

func (value *Value) SetFloat(v Float) {
    *(*Float)(pointer(value)) = v
    *(*uint32)(pointer(uintptr(pointer(value)) + 4)) = uint32(ValueTypeFloat) << 16
}

func (value *Value) GetBool() Bool {
    return *(*Bool)(pointer(value))
}

func (value *Value) SetBool(v Bool) {
    *(*Bool)(pointer(value)) = v
    *(*uint32)(pointer(uintptr(pointer(value)) + 4)) = uint32(ValueTypeBool) << 16
}

func (value Value) GetFunction() Function {
    if value.pointer == nil {
        return nil
    }

    return value.pointer.GetFunction()
}

func (value Value) GetObject() Object {
    val := value.GetInterface()
    if val == nil {
        return nil
    }
    return val.(Object)
}

func (value Value) String() string {
    switch value.GetType() {
    case ValueTypeInt:
        return strconv.FormatInt(int64(value.GetInt()), 10)
    case ValueTypeFloat:
        return strconv.FormatFloat(float64(value.GetFloat()), 'g', -1, 32)
    default:
        switch val := value.GetInterface().(type) {
        case String:
            return fmt.Sprintf("\"%v\"", string(val))
        default:
            return fmt.Sprintf("%v", val)
        }
    }
}

func (value Value) ToInt() Int {
    switch value.GetType() {
    case ValueTypeInt:
        return value.GetInt()
    case ValueTypeFloat:
        return Int(value.GetFloat())
    default:
        switch val := value.GetInterface().(type) {
        case Int64:
            return Int(val)
        case Float64:
            return Int(val)
        }
        return 0
    }
}

func (value Value) ToFloat() Float {
    switch value.GetType() {
    case ValueTypeFloat:
        return value.GetFloat()
    case ValueTypeInt:
        return Float(value.GetInt())
    case ValueTypeInterface:
        switch val := value.GetInterface().(type) {
        case Int64:
            return Float(val)
        case Float64:
            return Float(val)
        }
    }

    return Float(math.NaN())
}

func (value Value) ToInt64() Int64 {
    switch value.GetType() {
    case ValueTypeInterface:
        switch val := value.GetInterface().(type) {
        case Int64:
            return val
        case Float64:
            return Int64(val)
        }
    case ValueTypeInt:
        return Int64(value.GetInt())
    case ValueTypeFloat:
        return Int64(value.GetFloat())
    }

    return 0
}

func (value *Value) SetInt64(s Int64) {
    value.SetInterface(s)
}

func (value Value) ToFloat64() Float64 {
    switch value.GetType() {
    case ValueTypeInterface:
        switch val := value.GetInterface().(type) {
        case Float64:
            return val
        case Int64:
            return Float64(val)
        }
    case ValueTypeFloat:
        return Float64(value.GetFloat())
    case ValueTypeInt:
        return Float64(value.GetInt())
    }

    return Float64(math.NaN())
}

func (value *Value) SetFloat64(s Float64) {
    value.SetInterface(s)
}

func (value Value) ToBool() Bool {
    switch value.GetType() {
    case ValueTypeBool:
        return value.GetBool()
    case ValueTypeInt:
        return value.GetInt() != 0
    case ValueTypeFloat:
        return value.GetFloat() != 0
    case ValueTypeInterface:
        switch val := value.GetInterface().(type) {
        case Int64:
            return val != 0
        case Float64:
            return val != 0
        default:
            return true
        }
    }

    return false
}

func (value Value) ToString() string {
    switch value.GetType() {
    case ValueTypeInterface:
        return fmt.Sprint(value.GetInterface())
    }

    return ""
}

func (value Value) IsNull() bool {
    return value == NullValue || value.pointer == nil
}

func (value *Value) SetNull() {
    value.pointer = nil
}
