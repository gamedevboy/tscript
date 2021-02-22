package script

import (
	"unsafe"
)

type Bool bool

func (b Bool) MemorySize() int {
	return int(unsafe.Sizeof(b))
}

func (b Bool) Visit(memoryMap map[interface{}]int, f func(block MemoryBlock)) {
	f(b)
}

var _ MemoryBlock = Bool(false)

func (Bool) ScriptGet(fieldName string) interface{} {
	return NullValue
}

func (Bool) ScriptSet(string, interface{}) {
}

func (Bool) GetScriptTypeId() ScriptTypeId {
	return ScriptTypeBoolean
}
