package array

import (
	"tklibs/script"

	"tklibs/script/runtime"
	"tklibs/script/type/object"
)

type Component struct {
	elements []script.Value
	*object.Component
}

func (impl *Component) MemorySize() int {
	return 0
}

func (impl *Component) Visit(memoryMap map[interface{}]int, f func(block script.MemoryBlock)) {
	if _, ok := memoryMap[impl]; ok {
		return
	}

	memoryMap[impl] = impl.MemorySize()
	f(impl)

	impl.Component.Visit(memoryMap, f)

	for _, field := range impl.elements {
		if ms, ok := field.Get().(script.MemoryBlock); ok {
			ms.Visit(memoryMap, f)
		}
	}
}

func (impl *Component) First() interface{} {
	if len(impl.elements) == 0 {
		return script.Null
	}

	return impl.elements[0].Get()
}

func (impl *Component) Last() interface{} {
	if len(impl.elements) == 0 {
		return script.Null
	}

	return impl.elements[len(impl.elements)-1].Get()
}

func (impl *Component) Clear() {
	impl.elements = make([]script.Value, 0)
}

func (impl *Component) GetSlice() []script.Value {
	return impl.elements
}

var _ script.Array = &Component{}
var _ script.Object = &Component{}

func NewScriptArray(ctx interface{}, capSize int) *Component {
	ret := &Component{
		elements: make([]script.Value, 0, capSize),
	}
	ret.Component = object.NewScriptObject(ret, ctx, 0)
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
		return false
	}
	impl.elements = append(impl.elements[:index], impl.elements[index+1:]...)
	return true
}

func (impl *Component) Len() script.Int {
	return script.Int(len(impl.elements))
}
