package statement

type Case interface {
    GetValue() interface{}
    SetValue(value interface{})

    GetBlock() interface{}
    SetBlock(value interface{})
}
