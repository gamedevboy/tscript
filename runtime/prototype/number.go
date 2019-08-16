package prototype

import "tklibs/script/runtime"

type Number struct {
    prototype interface{}
}

func (impl *Number) GetNumberPrototype() interface{} {
    return impl.prototype
}

func NewNumberPrototype(ctx interface{}) *Number {
    ret := &Number{}
    ret.prototype = ctx.(runtime.ScriptContext).NewScriptObject(0)
    return ret
}
