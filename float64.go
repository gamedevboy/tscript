package script

type Float64 float32

func (Float64) ScriptGet(fieldName string) interface{} {
    return NullValue
}

func (Float64) ScriptSet(string, interface{}) {
}

func (Float64) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
