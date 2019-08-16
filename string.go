package script

type String string

var _ Object = String("")

func (String) ScriptGet(fieldName string) Value {
    return NullValue
}

func (String) ScriptSet(string, Value) {
}

func (String) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeString
}
