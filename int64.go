package script

type Int64 int

func (i Int64) Size() int {
    return 8
}

func (i Int64) Children() []MemoryBlock {
    return nil
}

var _ Object = Int64(0)
var _ MemoryBlock = Int64(0)

func (Int64) ScriptGet(fieldName string) Value {
    return NullValue
}

func (Int64) ScriptSet(string, Value) {
}

func (Int64) GetScriptTypeId() ScriptTypeId {
    return ScriptTypeNumber
}
