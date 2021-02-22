package script

import (
	"unsafe"
)

type Float float32

func (f Float) MemorySize() int {
	return int(unsafe.Sizeof(f))
}

func (f Float) Visit(memoryMap map[interface{}]int, f2 func(block MemoryBlock)) {
	f2(f)
}

var _ MemoryBlock = Float(0)

func (Float) ScriptGet(fieldName string) interface{} {
	return NullValue
}

func (Float) ScriptSet(string, interface{}) {
}

func (Float) GetScriptTypeId() ScriptTypeId {
	return ScriptTypeNumber
}
