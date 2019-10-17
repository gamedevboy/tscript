package frame

import (
    "sync"

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

var pool = sync.Pool{
    New: func() interface{} {
        return &Component{}
    },
}

func FreeStackFrame(c *Component) {
    pool.Put(c)
}

func NewStackFrame(owner, f interface{}) *Component {
    ret := pool.Get().(*Component)
    ret.ComponentType = script.MakeComponentType(owner)
    ret.function = f

    return ret
}
