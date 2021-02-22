package script

import (
	"unsafe"
)

type Int int32

func (i Int) Visit(memoryMap map[interface{}]int, f func(block MemoryBlock)) {
	f(i)
}

func (i Int) MemorySize() int {
	return int(unsafe.Sizeof(i))
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
