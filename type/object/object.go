package object

import (
    "tklibs/script"
    "tklibs/script/runtime"
    _ "tklibs/script/type/null"
)

type Component struct {
    script.ComponentType
    fields          []script.Value
    runtimeTypeInfo interface{}
    prototype       script.Value
}

func (impl *Component) GetPrototype() script.Value {
    return impl.prototype
}

func (impl *Component) SetPrototype(value script.Value) {
    impl.prototype = value
}

var _ script.Object = &Component{}
var _ runtime.Object = &Component{}

func (impl *Component) GetByIndex(index int) script.Value {
    return impl.fields[index]
}

func (impl *Component) SetByIndex(index int, value script.Value) {
    impl.fields[index] = value
}

func (*Component) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeObject
}

func (impl *Component) ScriptSet(fieldName string, value script.Value) {
    rt := impl.runtimeTypeInfo.(runtime.TypeInfo)
    index := rt.GetFieldIndexByName(fieldName)
    parent := rt.GetParent()
    if index > -1 {
        if value.IsNull() {
            if rt.GetName() == fieldName && parent != nil {
                impl.runtimeTypeInfo = parent
            } else {
                impl.runtimeTypeInfo = rt.RemoveChild(fieldName)
            }
            impl.fields = append(impl.fields[:index], impl.fields[index+1:]...)
        } else {
            impl.fields[index] = value
        }
    } else {
        if rt.GetName() == "-"+fieldName && parent != nil {
            impl.runtimeTypeInfo = parent
        } else {
            impl.runtimeTypeInfo = rt.AddChild(fieldName)
        }
        impl.fields = append(impl.fields, value)
    }
}

func (impl *Component) GetRuntimeTypeInfo() interface{} {
    return impl.runtimeTypeInfo
}

func (impl *Component) ScriptGet(fieldName string) script.Value {
    if fieldName == script.Prototype {
        return impl.prototype
    }

    v := impl.get(fieldName)
    if !v.IsNull() {
        return v
    }

    prototype := impl.prototype.GetInterface()
    if prototype == nil {
        return script.NullValue
    }

    return prototype.(script.Object).ScriptGet(fieldName)
}

func (impl *Component) get(s string) script.Value {
    index := impl.runtimeTypeInfo.(runtime.TypeInfo).GetFieldIndexByName(s)
    if index > -1 {
        return impl.fields[index]
    }

    return script.Value{}
}

func NewScriptObject(owner, ctx interface{}, fieldCap int) *Component {
    ret := &Component{
        ComponentType:   script.MakeComponentType(owner),
        runtimeTypeInfo: ctx.(runtime.ScriptContext).GetRootRuntimeType(),
        fields:          make([]script.Value, 0, fieldCap+1),
    }
    ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetObjectPrototype()))
    return ret
}

func NewScriptObjectWithRuntimePrototype(owner, runtimeType, prototype interface{}, fieldCap int) *Component {
    ret := &Component{
        ComponentType:   script.MakeComponentType(owner),
        runtimeTypeInfo: runtimeType,
        fields:          make([]script.Value, 0, fieldCap),
    }
    ret.SetPrototype(script.InterfaceToValue(prototype))
    return ret
}
