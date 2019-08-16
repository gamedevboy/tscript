package frame

import (
    "tklibs/script"
    "tklibs/script/runtime/stack"
)

type Component struct {
    script.ComponentType
    function interface{}
}

func (c *Component) GetFunction() interface{} {
    return c.function
}

var _ stack.Frame = &Component{}
var pool = make([]*Component, 0, 1)

func FreeStackFrame(c *Component) {
    pool = append(pool, c)
}

func NewStackFrame(owner, f interface{}) *Component {
    if len(pool) > 0 {
        ret := pool[len(pool)-1]
        pool = pool[:len(pool)-1]
        ret.function = f
        return ret
    }

    return &Component{
        ComponentType: script.MakeComponentType(owner),
        function:      f,
    }
}
