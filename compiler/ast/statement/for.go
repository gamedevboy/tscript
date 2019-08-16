package statement

type For interface {
    GetInit() interface{}
    SetInit(interface{})
    GetCondition() interface{}
    SetCondition(interface{})
    GetStep() interface{}
    SetStep(interface{})
    GetBody() interface{}
    SetBody(interface{})
}
