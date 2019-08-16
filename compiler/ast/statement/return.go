package statement

type Return interface {
    GetExpression() interface{}
    SetExpression(value interface{})
}
