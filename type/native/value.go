package native

import (
    "reflect"

    "tklibs/script"
)

type field struct {
    val    script.Value
    getter func(*Value, int) *script.Value
    setter func(*Value, int, *script.Value)
}

type Value struct {
    *script.ComponentType
    nativeType *Type
    value      reflect.Value
    fields     []field
}

func (v *Value) GetRuntimeTypeInfo() interface{} {
    return v.nativeType
}

func (v *Value) GetByIndex(index int) *script.Value {
    if index > -1 && index < len(v.fields) {
        return v.fields[index].getter(v, index)
    }

    ret := script.NullValue
    return &ret
}

func (v *Value) SetByIndex(index int, value *script.Value) {
    if index < 0 || index >= len(v.fields) {
        return
    }

    v.fields[index].setter(v, index, value)
}

func (v *Value) GetScriptTypeId() script.ScriptTypeId {
    return script.ScriptTypeObject
}

func (v *Value) ScriptSet(name string, value *script.Value) {
    index := v.nativeType.GetFieldIndexByName(name)
    if index > - 1 {
        v.SetByIndex(index, value)
    } else {
        f := field{val: *value}

        f.getter = func(*Value, int) *script.Value { return &f.val }
        f.setter = func(_ *Value, _ int, value *script.Value) { f.val = *value }

        v.fields = append(v.fields, f)

        v.nativeType = v.nativeType.AddChild(name).(*Type)
    }
}

func (v *Value) ScriptGet(name string) *script.Value {
    index := v.nativeType.GetFieldIndexByName(name)
    if index > -1 {
        return v.GetByIndex(index)
    }

    return script.NullValue
}

func NewValue(value reflect.Value) *Value {
    ret := &Value{}
    ret.ComponentType = script.NewComponentType(ret)
    ret.nativeType = GetType(value.Type())
    ret.value = value.Elem()
    ret.fields = make([]field, len(ret.nativeType.fields))

    copy(ret.fields, ret.nativeType.fields)

    for _, f := range ret.nativeType.methods {
        ret.ScriptSet(f.name, new(script.Value).SetInterface(&f))
    }

    return ret
}
