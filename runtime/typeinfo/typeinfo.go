package typeinfo

import (
    "tklibs/script"
    "tklibs/script/runtime"
)

type Component struct {
    script.ComponentType
    parent    interface{}
    children  map[string]interface{}
    fieldName string
    fields    []string
}

var _ runtime.TypeInfo = &Component{}

func NewTypeComponent(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        fields:        make([]string, 0),
        children:      make(map[string]interface{}),
    }
}

func (impl *Component) GetParent() interface{} {
    return impl.parent
}

func (impl *Component) GetName() string {
    return impl.fieldName
}

func (impl *Component) GetFieldIndexByName(fieldName string) int {
    for index, fn := range impl.fields {
        if fn == fieldName {
            return index
        }
    }

    return -1
}

func (impl *Component) GetFieldNames() []string {
    return impl.fields
}

func (impl *Component) GetFieldValueByIndex(index int) interface{} {
    return impl.fields
}

func (impl *Component) AddChild(fieldName string) interface{} {
    if v, ok := impl.children[fieldName]; ok {
        return v
    }

    newTypeInfo := &struct {
        *Component
    }{}
    newTypeInfo.Component = NewTypeComponent(newTypeInfo)
    newTypeInfo.fieldName = fieldName
    newTypeInfo.parent = impl.GetOwner()

    newTypeInfo.fields = make([]string, len(impl.fields), len(impl.fields)+1)
    copy(newTypeInfo.fields, impl.fields)
    newTypeInfo.fields = append(newTypeInfo.fields, fieldName)

    impl.children[fieldName] = newTypeInfo

    return newTypeInfo
}

func (impl *Component) RemoveChild(fieldName string) interface{} {
    fieldName = "-" + fieldName
    if v, ok := impl.children[fieldName]; ok {
        return v
    }

    newTypeInfo := &struct {
        *Component
    }{}
    newTypeInfo.Component = NewTypeComponent(newTypeInfo)
    newTypeInfo.fieldName = fieldName
    newTypeInfo.parent = impl.GetOwner()

    for index, fn := range impl.fields {
        if fn == fieldName {
            newTypeInfo.fields = make([]string, len(impl.fields)-1, len(impl.fields)-1)
            copy(newTypeInfo.fields[:index], impl.fields[:index])
            copy(newTypeInfo.fields[index:], impl.fields[index+1:])
            break
        }
    }

    impl.children[fieldName] = newTypeInfo

    return newTypeInfo
}
