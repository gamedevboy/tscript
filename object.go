package script

type Object interface {
    GetScriptTypeId() ScriptTypeId
    ScriptSet(string, Value)
    ScriptGet(string) Value
}
