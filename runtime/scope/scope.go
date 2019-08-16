package scope

import (
    "container/list"

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
var pool = make([]*Component, 0, 1)

func FreeScope(c *Component) {
    pool = append(pool, c)
}

func NewScope(owner, function interface{}, argList, localVarList []script.Value) *Component {
    if len(pool) > 0 {
        ret := pool[len(pool)-1]
        pool = pool[:len(pool)-1]
        ret.function = function
        ret.argList = argList
        ret.localVarList = localVarList
        ret.refListMap = make(map[*script.Value]*list.List)
        return ret
    }
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        function:      function,
        argList:       argList,
        localVarList:  localVarList,
        refListMap:    make(map[*script.Value]*list.List),
    }
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
