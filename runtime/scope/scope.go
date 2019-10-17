package scope

import (
    "container/list"
    "sync"

    "tklibs/script"
    "tklibs/script/runtime"
)

type Component struct {
    script.ComponentType
    function     interface{}
    argList      []script.Value
    localVarList []script.Value
    refListMap   map[*script.Value]*list.List
}

func (impl *Component) KeepRefs() {
    for k, v := range impl.refListMap {
        r := new(script.Value)
        *r = *k
        for it := v.Front(); it != nil; it = it.Next() {
            *it.Value.(**script.Value) = r
        }
    }
}

func (impl *Component) AddToRefList(value *script.Value, valuePtr **script.Value) {
    *valuePtr = value
    if l, ok := impl.refListMap[value]; ok {
        l.PushBack(valuePtr)
    } else {
        l = list.New()
        l.PushBack(valuePtr)
        impl.refListMap[value] = l
    }
}

var _ runtime.Scope = &Component{}
var pool = sync.Pool{New: func() interface{} {
    return &Component{}
}}

func FreeScope(c *Component) {
    pool.Put(c)
}

func NewScope(owner, function interface{}, argList, localVarList []script.Value) *Component {
    var comp = pool.Get().(*Component)

    comp.ComponentType = script.MakeComponentType(owner)
    comp.function = function
    comp.argList = argList
    comp.localVarList = localVarList
    comp.refListMap = make(map[*script.Value]*list.List)

    return comp
}

func (impl *Component) GetFunction() interface{} {
    return impl.function
}

func (impl *Component) GetArgList() []script.Value {
    return impl.argList
}

func (impl *Component) GetLocalVarList() []script.Value {
    return impl.localVarList
}
