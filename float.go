package script

type Float float32

func (Float) ScriptGet(fieldName string) interface{} {
    return NullValue
}

func (Float) ScriptSet(string, interface{}) {
}

func (Float) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
