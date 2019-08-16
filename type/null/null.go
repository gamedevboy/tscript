package null

import (
    "fmt"

    "tklibs/script"
)

type scriptNull struct {
    script.ComponentType
}

var _ script.Object = &scriptNull{}

func (*scriptNull) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeNull
}

func (*scriptNull) ScriptGet(name string) script.Value {
    panic(fmt.Errorf("Can not get '%v' from null ", name))
}

func (*scriptNull) ScriptSet(name string, val script.Value) {
    panic(fmt.Errorf("Can not set '%v' to null ", name))
}

func (*scriptNull) String() string {
    return "null"
}

var _ script.Object = &scriptNull{}

func init() {
    null := &struct {
        *scriptNull
    }{}
    null.scriptNull = &scriptNull{script.MakeComponentType(null)}

    script.Null = null

    script.NullValue.SetInterface(null)
}
