package runtime

import (
    "tklibs/script"
)

type ScriptContext interface {
    PushFrame(frame interface{})
    PopFrame() interface{}

    PushScope(interface{})
    PopScope() interface{}

    GetCurrentFrame() interface{}

    GetRootRuntimeType() interface{}

    NewScriptObject(fieldCap int) interface{}
    NewScriptArray(sizeCap int) interface{}
    NewScriptMap(sizeCap int) interface{}

    GetFunctionPrototype() interface{}
    GetNumberPrototype() interface{}
    GetObjectPrototype() interface{}
    GetStringPrototype() interface{}
    GetBoolPrototype() interface{}
    GetArrayPrototype() interface{}
    GetMapPrototype() interface{}

    GetAssembly() interface{}

    Run() interface{}
    Reload(interface{}) interface{}

    GetRefByName(name string, valuePtr **script.Value)
    GetRegisters() []script.Value
    PushRegisters(regStart script.Int, length int) []script.Value
    PopRegisters()
}
