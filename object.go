package script

type Object interface {
	GetScriptTypeId() ScriptTypeId
	ScriptSet(string, Value)
	ScriptGet(string) Value
}

type MemoryBlock interface {
	MemorySize() int
	Visit(memoryMap map[interface{}]int, f func(block MemoryBlock))
}
