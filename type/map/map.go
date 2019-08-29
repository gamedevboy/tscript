package _map

import (
    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/type/object"
)

type Component struct {
    script.ComponentType
    *object.Component

    values map[script.Value]script.Value
}

func (impl *Component) Set(key, value interface{}) {
    k, v := script.Value{}, script.Value{}
    k.Set(key)
    v.Set(v)
    impl.values[k] = v
}

func (impl *Component) Get(key interface{}) interface{} {
    k := script.Value{}
    k.Set(key)
    if ret, ok := impl.values[k]; ok {
        return ret.Get()
    }

    return script.Null
}

var _ script.Object = &Component{}
var _ script.Map = &Component{}

func (impl *Component) Delete(key script.Value) {
    if key.GetType() == script.ValueTypeInterface {
        if keyVal, ok := key.GetInterface().(script.String); ok {
            impl.ScriptSet(string(keyVal), script.Value{})
            return
        }
    }

    delete(impl.values, key)
}

func (impl *Component) Len() script.Int {
    return script.Int(len(impl.values) + len(impl.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldNames()))
}

func (impl *Component) ContainsKey(value script.Value) script.Bool {
    _, ok := impl.values[value]
    return script.Bool(ok)
}

func (impl *Component) Foreach(f func(key, value interface{}) bool) script.Int {
    var i script.Int
    for key, value := range impl.values {
        if !f(key.Get(), value.Get()) {
            break
        }
        i++
    }

    typeInfo := impl.GetRuntimeTypeInfo().(runtime.TypeInfo)
    for i, key := range typeInfo.GetFieldNames() {
        if !f(script.String(key), impl.GetByIndex(i).Get()) {
            break
        }
    }

    return i
}

var _ script.Map = &Component{}

func (impl *Component) SetValue(key, value script.Value) {
    switch key.GetType() {
    case script.ValueTypeInterface:
        switch k := key.GetInterface().(type) {
        case script.String:
            impl.ScriptSet(string(k), value)
            return
        }
        fallthrough
    default:
        impl.values[key] = value
    }
}

func (impl *Component) GetValue(key script.Value) script.Value {
    switch key.GetType() {
    case script.ValueTypeInterface:
        switch k := key.GetInterface().(type) {
        case script.String:
            return impl.ScriptGet(string(k))
        }
        fallthrough
    default:
        if value, ok := impl.values[key]; ok {
            return value
        }

        return script.NullValue
    }
}

func (*Component) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeMap
}

func NewScriptMap(owner, ctx interface{}, capSize int) *Component {
    ret := &Component{
        ComponentType: script.MakeComponentType(owner),
        Component:     object.NewScriptObject(owner, ctx, 0),
        values:        make(map[script.Value]script.Value, capSize),
    }

    ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetMapPrototype()))
    return ret
}
