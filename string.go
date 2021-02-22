package script

type String string

func (s String) Visit(memoryMap map[interface{}]int, f func(block MemoryBlock)) {
	f(s)
}

var _ Object = String("")
var _ MemoryBlock = String("")

func (String) ScriptGet(fieldName string) Value {
	return NullValue
}

func (s String) MemorySize() int {
	return len(s)
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
