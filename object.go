package script

type Object interface {
    GetScriptTypeId() ScriptTypeId
    ScriptSet(string, Value)
    ScriptGet(string) Value
}

type MemoryBlock interface {
    Size() int
    Children() []MemoryBlock
}
