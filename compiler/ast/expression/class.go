package expression

import "container/list"

type Class interface {
    Call
    SetName(string)
    GetName() string
    SetParent(interface{})
    GetParent() interface{}
    GetMethods() *list.List
    AddMethod(f interface{})
    FinishMethod()
}
