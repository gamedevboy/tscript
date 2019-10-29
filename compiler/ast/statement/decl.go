package statement

type Decl interface {
    GetName() string
    SetName(string)

    GetExpression() interface{}
    SetExpression(interface{})

    IsGlobal() bool
    SetGlobal(bool)
}

