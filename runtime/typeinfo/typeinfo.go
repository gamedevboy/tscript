package typeinfo

import (
	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/runtime/util"
)

type Component struct {
	script.ComponentType
	parent     interface{}
	children   map[*string]interface{}
	fieldName  *string
	fields     []*string
	stringPool util.StringPool
}

var _ runtime.TypeInfo = &Component{}

func NewTypeComponent(owner interface{}, stringPool util.StringPool) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
		fields:        make([]*string, 0),
		children:      make(map[*string]interface{}),
		stringPool:    stringPool,
		fieldName:     nil,
	}
}

func (impl *Component) GetParent() interface{} {
	return impl.parent
}

func (impl *Component) GetName() string {
	if impl.fieldName == nil {
		return ""
	}

	return *impl.fieldName
}

func (impl *Component) GetFieldIndexByName(fieldName string) int {
	for index, fn := range impl.fields {
		if *fn == fieldName {
			return index
		}
	}

	return -1
}

func (impl *Component) GetFieldNames() []*string {
	return impl.fields
}

func (impl *Component) GetFieldValueByIndex(index int) interface{} {
	return impl.fields
}

func (impl *Component) AddChild(fieldName string) interface{} {
	fn := impl.stringPool.Insert(fieldName)
	if v, ok := impl.children[fn]; ok {
		return v
	}

	newTypeInfo := &struct {
		*Component
	}{}
	newTypeInfo.Component = NewTypeComponent(newTypeInfo, impl.stringPool)
	newTypeInfo.fieldName = fn
	newTypeInfo.parent = impl.GetOwner()

	newTypeInfo.fields = make([]*string, len(impl.fields), len(impl.fields)+1)
	copy(newTypeInfo.fields, impl.fields)
	newTypeInfo.fields = append(newTypeInfo.fields, fn)

	impl.children[fn] = newTypeInfo

	return newTypeInfo
}

func (impl *Component) RemoveChild(fieldName string) interface{} {
	newFieldName := "-" + fieldName
	nfn := impl.stringPool.Insert(newFieldName)
	if v, ok := impl.children[nfn]; ok {
		return v
	}

	newTypeInfo := &struct {
		*Component
	}{}
	newTypeInfo.Component = NewTypeComponent(newTypeInfo, impl.stringPool)
	newTypeInfo.fieldName = nfn
	newTypeInfo.parent = impl.GetOwner()

	for index, fn := range impl.fields {
		if fn == nfn {
			newTypeInfo.fields = make([]*string, len(impl.fields)-1, len(impl.fields)-1)
			copy(newTypeInfo.fields[:index], impl.fields[:index])
			copy(newTypeInfo.fields[index:], impl.fields[index+1:])
			break
		}
	}

	impl.children[nfn] = newTypeInfo

	return newTypeInfo
}
