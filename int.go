package script

type Int int32

func (Int) ScriptGet(fieldName string) Value {
    return NullValue
}

func (Int) ScriptSet(string, Value) {
}

func (Int) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
