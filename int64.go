package script

type Int64 int

func (Int64) ScriptGet(fieldName string) Value {
    return NullValue
}

func (Int64) ScriptSet(string, Value) {
}

func (Int64) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
