package script

type String string

var _ Object = String("")
var _ MemoryBlock = String("")

func (String) ScriptGet(fieldName string) Value {
    return NullValue
}

func (s String) Size() int {
    return len(s)
}

func (s String) Children() []MemoryBlock {
    return nil
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
