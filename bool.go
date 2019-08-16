package script

type Bool bool

func (Bool) ScriptGet(fieldName string) interface{} {
    return NullValue
}

func (Bool) ScriptSet(string, interface{}) {
}

func (Bool) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeBoolean
}
