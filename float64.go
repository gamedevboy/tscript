package script

import (
	"unsafe"
)

type Float64 float32

func (f Float64) MemorySize() int {
	return int(unsafe.Sizeof(f))
}

func (f Float64) Visit(memoryMap map[interface{}]int, f2 func(block MemoryBlock)) {
	f2(f)
}

var _ MemoryBlock = Float64(0)

func (Float64) ScriptGet(fieldName string) interface{} {
	return NullValue
}

func (Float64) ScriptSet(string, interface{}) {
}

func (Float64) GetScriptTypeId() ScriptTypeId {
	return ScriptTypeNumber
}
