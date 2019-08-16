package script

type Function interface {
    Invoke(interface{}, ...interface{}) interface{}
    SetThis(Value)
    GetThis() Value

    GetRuntimeFunction() interface{}

    GetRefList() []*Value

    GetFieldByMemberIndex(obj interface{}, index Int) Value
    SetFieldByMemberIndex(obj interface{}, index Int, value Value)
}
