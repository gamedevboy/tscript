package script

import (
	"unsafe"
)

type Int64 int

func (i Int64) Visit(memoryMap map[interface{}]int, f func(block MemoryBlock)) {
	f(i)
}

func (i Int64) MemorySize() int {
	return int(unsafe.Sizeof(i))
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
