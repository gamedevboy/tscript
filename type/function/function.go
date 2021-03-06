package function

import (
	"container/list"
	"fmt"
	"reflect"

	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/runtime/native"
	"tklibs/script/runtime/runtime_t"
	"tklibs/script/type/object"
)

var _ script.Object = &Component{}
var _ script.MemoryBlock = &Component{}

type Component struct {
	refs          []*script.Value
	memberCaches  []*list.List
	refNames      []string
	scriptContext interface{}

	script.ComponentType
	*object.Component
	this script.Value

	scriptRuntimeFunction runtime_t.Function
	nativeRuntimeFunction runtime_t.NativeFunction
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

	for _, field := range impl.refs {
		if ms, ok := field.Get().(script.MemoryBlock); ok {
			ms.Visit(memoryMap, f)
		}
	}
}

func (impl *Component) GetContext() interface{} {
	return impl.scriptContext
}

func (impl *Component) IsScriptFunction() bool {
	return impl.scriptRuntimeFunction != nil
}

func (impl *Component) GetScriptRuntimeFunction() runtime_t.Function {
	return impl.scriptRuntimeFunction
}

func (impl *Component) GetNativeRuntimeFunction() runtime_t.NativeFunction {
	return impl.nativeRuntimeFunction
}

func (impl *Component) SetThis(this script.Value) {
	impl.this = this
}

func (impl *Component) GetThis() script.Value {
	return impl.this
}

var _ script.Function = &Component{}

type fieldCache struct {
	cacheType runtime.TypeInfo
	offset    int
}

var invalidFieldCache = &fieldCache{
	offset: -1,
}

func (*Component) GetScriptTypeId() script.ScriptTypeId {
	return script.ScriptTypeFunction
}

func (impl *Component) GetRefList() []*script.Value {
	return impl.refs
}

func (impl *Component) getFieldCache(obj interface{}, index script.Int) *fieldCache {
	switch runtimeObj := obj.(type) {
	case runtime.Object:
		for it := impl.memberCaches[index].Front(); it != nil; it = it.Next() {
			if it.Value.(*fieldCache).cacheType == runtimeObj.GetRuntimeTypeInfo() {
				return it.Value.(*fieldCache)
			}
		}

		rt := runtimeObj.GetRuntimeTypeInfo().(runtime.TypeInfo)
		offset := rt.GetFieldIndexByName(impl.getMemberNames()[int(index)])
		mc := &fieldCache{rt, offset}

		impl.memberCaches[index].PushBack(mc)

		return mc
	default:
		return invalidFieldCache
	}
}

func (impl *Component) Invoke(context, this interface{}, args ...interface{}) interface{} {
	return context.(runtime.ScriptInterpreter).InvokeFunction(impl.GetOwner(), context, this, args...)
}

func (impl *Component) New(context interface{}, args ...interface{}) interface{} {
	if impl.scriptRuntimeFunction != nil {
		for index, value := range impl.refs {
			if value.GetType() == script.ValueTypeInterface && value.GetInterface() == nil {
				panic(fmt.Sprintf("Can not find ref value for '%v'", impl.scriptRuntimeFunction.GetRefVars()[index]))
			}
		}
	}

	return impl.scriptContext.(runtime.ScriptInterpreter).InvokeNew(impl.GetOwner(), context, args...)
}

func (impl *Component) Init() {
	if impl.scriptRuntimeFunction != nil {
		scriptContext := impl.scriptContext.(runtime.ScriptContext)

		impl.refNames = impl.scriptRuntimeFunction.GetRefVars()
		refLength := len(impl.refNames)
		impl.refs = make([]*script.Value, refLength)

		for i := 0; i < refLength; i++ {
			scriptContext.GetRefByName(impl.refNames[i], &impl.refs[i])
		}

		impl.memberCaches = make([]*list.List, len(impl.scriptRuntimeFunction.GetMembers()))
		for i := range impl.memberCaches {
			impl.memberCaches[i] = list.New()
		}
	}
}

func (impl *Component) getMemberNames() []string {
	return impl.scriptRuntimeFunction.GetMembers()
}

func NewScriptFunction(owner, runtimeFunction, ctx interface{}) *Component {
	ret := &Component{
		ComponentType: script.MakeComponentType(owner),
		scriptContext: ctx,
		Component:     object.NewScriptObject(owner, ctx, 0),
	}

	switch _func := runtimeFunction.(type) {
	case runtime_t.Function:
		ret.scriptRuntimeFunction = _func
	case runtime_t.NativeFunction:
		ret.nativeRuntimeFunction = _func
	}

	ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetFunctionPrototype()))

	return ret
}

func NativeFunctionToValue(fn native.FunctionType, context interface{}) script.Value {
	return script.InterfaceToValue(NewNativeFunction(fn, context))
}

func NewNativeFunction(f, ctx interface{}) interface{} {
	nativeFunction := &struct {
		*Component
	}{}

	switch _func := f.(type) {
	case native.FunctionType:
		nativeFunction.Component = NewScriptFunction(nativeFunction, _func, ctx)
	case reflect.Value:
		nativeFunction.Component = NewScriptFunction(nativeFunction, native.ReflectFunction(_func), ctx)
	default:
		panic("")
	}

	return nativeFunction
}
