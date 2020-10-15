package function

import (
	"container/list"
	"fmt"
	"reflect"
	rt "runtime"
	"unsafe"

	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/runtime/native"
	"tklibs/script/runtime/runtime_t"
	"tklibs/script/type/object"
)

type ref struct {
	pointer uintptr
	fn      runtime_t.Function
}

func (comp *ref) Dispose() {
	comp.fn.UnregisterFunction(comp.pointer)
}

type Component struct {
	refs          []*script.Value
	memberCaches  []*list.List
	refNames      []string
	scriptContext interface{}

	script.ComponentType
	*object.Component
	this     script.Value
	_compRef *ref

	scriptRuntimeFunction runtime_t.Function
	nativeRuntimeFunction runtime_t.NativeFunction
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

func (impl *Component) Reload() {
	if impl.scriptRuntimeFunction != nil {
		scriptContext := impl.scriptContext.(runtime.ScriptContext)

		impl.memberCaches = make([]*list.List, len(impl.scriptRuntimeFunction.GetMembers()))
		for i := range impl.memberCaches {
			impl.memberCaches[i] = list.New()
		}

		refs := impl.refs
		refNames := impl.refNames

		impl.refNames = impl.scriptRuntimeFunction.GetRefVars()
		refLength := len(impl.refNames)
		impl.refs = make([]*script.Value, refLength)

		for i := 0; i < refLength; i++ {
			found := false

			for k := 0; k < len(refs); k++ {
				if refNames[k] == impl.refNames[i] {
					impl.refs[i] = refs[k]
					found = true
					break
				}
			}

			if !found {
				scriptContext.GetRefByName(impl.refNames[i], &impl.refs[i])
			}
		}
	}
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

func (impl *Component) GetFieldByMemberIndex(obj interface{}, index script.Int) script.Value {
	switch target := obj.(type) {
	case script.Int, script.Int64:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetNumberPrototype(), index)
	case script.Float, script.Float64:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetNumberPrototype(), index)
	case script.Bool:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetBoolPrototype(), index)
	case script.String:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetStringPrototype(), index)
	case script.Map:
		key := script.String(impl.getMemberNames()[index])
		if target.ContainsKey(key) {
			return script.ToValue(target.Get(key))
		}
		return obj.(script.Object).ScriptGet(impl.getMemberNames()[index])
	case script.Object:
		offset := impl.getFieldCache(obj, index).offset
		if offset > -1 {
			return *obj.(runtime.Object).GetByIndex(offset)
		}

		return target.ScriptGet(impl.getMemberNames()[index])
	}

	return script.NullValue
}

func (impl *Component) SetFieldByMemberIndex(obj interface{}, index script.Int, value script.Value) {
	switch target := obj.(type) {
	case script.Map:
		target.Set(script.String(impl.getMemberNames()[index]), value.Get())
	default:
		offset := impl.getFieldCache(obj, index).offset
		if offset > -1 {
			obj.(runtime.Object).SetByIndex(offset, value)
		} else {
			obj.(script.Object).ScriptSet(impl.getMemberNames()[index], value)
		}
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
		ret._compRef = &ref{
			pointer: ^uintptr(unsafe.Pointer(ret)),
			fn:      _func,
		}

		_func.RegisterFunction(ret._compRef.pointer)
		rt.SetFinalizer(ret._compRef, (*ref).Dispose)
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
