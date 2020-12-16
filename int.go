package script

type Int int32

func (i Int) Size() int {
    return 4
}

func (i Int) Children() []MemoryBlock {
    return nil
}

var _ Object = Int(0)
var _ MemoryBlock = Int(0)

func (Int) ScriptGet(fieldName string) Value {
    return NullValue
}

func (Int) ScriptSet(string, Value) {
}

func (Int) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
