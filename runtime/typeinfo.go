package runtime

type TypeInfo interface {
    AddChild(fieldName string) interface{}
    RemoveChild(fieldName string) interface{}
    GetName() string
    GetParent() interface{}
    GetFieldIndexByName(fieldName string) int
    GetFieldNames() []*string
}
