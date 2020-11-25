package typeinfo

import (
	"tklibs/script/runtime"
)

type Component struct {
	fields        []string
	scriptContext runtime.ScriptContext
	parent        runtime.TypeInfo
	children      map[string]runtime.TypeInfo
}

func (impl *Component) GetContext() runtime.ScriptContext {
    return impl.scriptContext
}

var _ runtime.TypeInfo = &Component{}

func NewTypeComponent(scriptContext interface{}) *Component {
	return &Component{
		fields:        make([]string, 0),
		children:      make(map[string]runtime.TypeInfo),
		scriptContext: scriptContext.(runtime.ScriptContext),
	}
}

func (impl *Component) GetParent() runtime.TypeInfo {
	return impl.parent
}

func (impl *Component) GetName() string {
	fl := len(impl.fields)

	if fl > 0 {
		return impl.fields[fl-1]
	}

	return ""
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

func (impl *Component) AddChild(fieldName string) runtime.TypeInfo {
	fn := impl.scriptContext.GetStringPool().Insert(fieldName)
	if v, ok := impl.children[fn]; ok {
		return v
	}

	newTypeInfo := NewTypeComponent(impl.scriptContext)
	newTypeInfo.parent = impl

	newTypeInfo.fields = make([]string, len(impl.fields), len(impl.fields)+1)
	copy(newTypeInfo.fields, impl.fields)
	newTypeInfo.fields = append(newTypeInfo.fields, fn)

	impl.children[fn] = newTypeInfo

	return newTypeInfo
}

func (impl *Component) RemoveChild(fieldName string) runtime.TypeInfo {
	newFieldName := "-" + fieldName
	nfn := impl.scriptContext.GetStringPool().Insert(newFieldName)
	if v, ok := impl.children[nfn]; ok {
		return v
	}

	newTypeInfo := NewTypeComponent(impl.scriptContext)
	newTypeInfo.parent = impl

	for index, fn := range impl.fields {
		if fn == nfn {
			newTypeInfo.fields = make([]string, len(impl.fields)-1, len(impl.fields)-1)
			copy(newTypeInfo.fields[:index], impl.fields[:index])
			copy(newTypeInfo.fields[index:], impl.fields[index+1:])
			break
		}
	}

	impl.children[nfn] = newTypeInfo

	return newTypeInfo
}
