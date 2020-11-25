package runtime

type TypeInfo interface {
    AddChild(fieldName string) TypeInfo
    RemoveChild(fieldName string) TypeInfo
    GetName() string
    GetParent() TypeInfo
    GetFieldIndexByName(fieldName string) int
    GetFieldNames() []string
    GetContext() ScriptContext
}
