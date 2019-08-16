package library

type RuntimeLibrary interface {
    GetName() string
    SetScriptContext(context interface{})
}
