package script

type ScriptTypeId uint8

const (
    ScriptTypeNumber ScriptTypeId = iota
    ScriptTypeBoolean
    ScriptTypeString
    ScriptTypeObject
    ScriptTypeFunction
    ScriptTypeArray
    ScriptTypeMap
    ScriptTypeNull
)

const Prototype = "prototype"
