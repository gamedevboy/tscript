package script

import "unsafe"

const (
    InterfaceTypeAny uint64 = iota
    InterfaceTypeInt64
    InterfaceTypeFloat64
    InterfaceTypeString
    InterfaceTypeArray
    InterfaceTypeMap
    InterfaceTypeFunction
    InterfaceTypeObject
)

type Interface struct {
    value *interface{}
    iType uint64
}

func (i *Interface) Set(v interface{}) {
    switch value := v.(type) {
    case Function:
        *(**Function)(unsafe.Pointer(i)) = &value
        i.iType = InterfaceTypeFunction
    default:
        i.value = &v
        i.iType = InterfaceTypeAny
    }
}

func (i *Interface) GetFunction() Function {
    return **(**Function)(unsafe.Pointer(i))
}

func (i *Interface) Get() interface{} {
    switch i.iType {
    case InterfaceTypeFunction:
        return i.GetFunction()
    default:
        return *i.value
    }
}

func (i *Interface) GetType() uint64 {
    return i.iType
}

func (i *Interface) GetObject() Object {
    return **(**Object)(unsafe.Pointer(i))
}
