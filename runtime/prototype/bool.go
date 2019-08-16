package prototype

import "tklibs/script/runtime"

type Bool struct {
    prototype interface{}
}

func (impl *Bool) GetBoolPrototype() interface{} {
    return impl.prototype
}

func NewBoolPrototype(ctx interface{}) *Bool {
    ret := &Bool{}
    ret.prototype = ctx.(runtime.ScriptContext).NewScriptObject(0)
    return ret
}
