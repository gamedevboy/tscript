package array

import (
    "tklibs/script"

    "tklibs/script/runtime"
    "tklibs/script/type/object"
)

type Component struct {
    script.ComponentType
    *object.Component
    elements []script.Value
}

func (impl *Component) Clear() {
    impl.elements = make([]script.Value, 0)
}

func (impl *Component) GetSlice() interface{} {
    return impl.elements
}

var _ script.Array = &Component{}
var _ script.Object = &Component{}

func NewScriptArray(owner, ctx interface{}, capSize int) *Component {
    ret := &Component{
        ComponentType: script.MakeComponentType(owner),
        Component:     object.NewScriptObject(owner, ctx, 0),
        elements:      make([]script.Value, 0, capSize),
    }
    ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetArrayPrototype()))
    return ret
}

func (impl *Component) Push(args ...script.Value) interface{} {
    impl.elements = append(impl.elements, args...)
    return impl.GetOwner()
}

func (impl *Component) Pop() interface{} {
    lastIndex := len(impl.elements) - 1
    ret := impl.elements[lastIndex]
    impl.elements = impl.elements[:lastIndex]
    return ret.Get()
}

func (impl *Component) Shift() interface{} {
    if len(impl.elements) < 1 {
        return script.Null
    }

    ret := impl.elements[0]
    impl.elements = impl.elements[1:]
    return ret.Get()
}

func (impl *Component) Unshift(args ...script.Value) interface{} {
    impl.elements = append(args, impl.elements...)
    return impl.GetOwner()
}

func (*Component) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeArray
}

func (impl *Component) GetElement(index script.Int) script.Value {
    if index < 0 || int(index) >= len(impl.elements) {
        return script.NullValue
    }
    return impl.elements[index]
}

func (impl *Component) SetElement(index script.Int, value script.Value) {
    size := len(impl.elements)
    if int(index) >= size {
        for i := size - 1; i < int(index); i++ {
            impl.elements = append(impl.elements, script.NullValue)
        }
    }

    impl.elements[index] = value
}

func (impl *Component) RemoveAt(index script.Int) script.Bool {
    if index < 0 || index >= script.Int(len(impl.elements)) {
        return script.Bool(false)
    }
    impl.elements = append(impl.elements[:index], impl.elements[index+1:]...)
    return true
}

func (impl *Component) Len() script.Int {
    return script.Int(len(impl.elements))
}
