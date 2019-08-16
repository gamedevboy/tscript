package statement

type If interface {
    GetCondition() interface{}
    SetCondition(interface{})
    GetBody() interface{}
    SetBody(interface{})
    GetElseBody() interface{}
    SetElseBody(interface{})
}
