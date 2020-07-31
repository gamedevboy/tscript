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

func (s String) ToValue() Value {
    v := Value{}
    v.SetInterface(s)
    return v
}
