package _map

import (
	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/type/object"
)

type Component struct {
	*object.Component
	values map[interface{}]interface{}
}

func (impl *Component) Set(key, value interface{}) {
	impl.values[key] = value
}

func (impl *Component) Get(key interface{}) interface{} {
	if value, ok := impl.values[key]; ok {
		return value
	}
	return script.Null
}

var _ script.Object = &Component{}
var _ script.Map = &Component{}

func (impl *Component) Delete(key interface{}) {
	delete(impl.values, key)
}

func (impl *Component) Len() script.Int {
	return script.Int(len(impl.values))
}

func (impl *Component) ContainsKey(key interface{}) script.Bool {
	_, ok := impl.values[key]
	return script.Bool(ok)
}

func (impl *Component) Foreach(f func(key, value interface{}) bool) script.Int {
	var i script.Int
	for key, value := range impl.values {
		if !f(key, value) {
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

func (*Component) GetScriptTypeId() script.ScriptTypeId {
	return script.ScriptTypeMap
}

func NewScriptMap(ctx interface{}, capSize int) *Component {
	ret := &Component{
		values:        make(map[interface{}]interface{}, capSize),
	}

	ret.Component = object.NewScriptObject(ret, ctx, 0)
	ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetMapPrototype()))

	return ret
}
