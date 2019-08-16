package statement

type While interface {
    GetCondition() interface{}
    SetCondition(interface{})
    GetBody() interface{}
    SetBody(interface{})
}
