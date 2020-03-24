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

type ComopnentRef struct {
	pointer uintptr
	f       runtime_t.Function
}

func (comp *ComopnentRef) Dispose() {
	if comp.pointer != 0 {
		comp.f.UnregisterFunction(comp.pointer)
	}
}

type Component struct {
	script.ComponentType
	*object.Component

	_compRef *ComopnentRef

	scriptContext   interface{}
	runtimeFunction interface{}

	scriptRuntimeFunction runtime_t.Function
	nativeRuntimeFunction runtime_t.NativeFunction

	scriptFunction bool

	refs         []*script.Value
	memberCaches []*list.List
	this         script.Value
	refNames     []string
}

func (impl *Component) GetContext() interface{} {
	return impl.scriptContext
}

func (impl *Component) IsScriptFunction() bool {
	return impl.scriptFunction
}

func (impl *Component) GetScriptRuntimeFunction() runtime_t.Function {
	return impl.scriptRuntimeFunction
}

func (impl *Component) GetNativeRuntimeFunction() runtime_t.NativeFunction {
	return impl.nativeRuntimeFunction
}

func (impl *Component) Reload() {
	if runtimeFunction, ok := impl.runtimeFunction.(runtime_t.Function); ok {
		scriptContext := impl.scriptContext.(runtime.ScriptContext)

		impl.memberCaches = make([]*list.List, len(runtimeFunction.GetMembers()))
		for i := range impl.memberCaches {
			impl.memberCaches[i] = list.New()
		}

		refs := impl.refs
		refNames := impl.refNames

		impl.refNames = runtimeFunction.GetRefVars()
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
	case script.Int:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetNumberPrototype(), index)
	case script.Float:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetNumberPrototype(), index)
	case script.Bool:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetBoolPrototype(), index)
	case script.String:
		return impl.GetFieldByMemberIndex(impl.scriptContext.(runtime.ScriptContext).GetStringPrototype(), index)
	case script.Map:
		key := script.String(impl.getMemberNames()[index])
		if target.ContainsKey(key) {
			val := script.Value{}
			val.Set(target.Get(key))
			return val
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

func (impl *Component) GetRuntimeFunction() interface{} {
	return impl.runtimeFunction
}

func (impl *Component) Invoke(this interface{}, args ...interface{}) interface{} {
	return impl.scriptContext.(runtime.ScriptInterpreter).InvokeFunction(impl.GetOwner(), this, args...)
}

func (impl *Component) New(args ...interface{}) interface{} {
	switch runtimeFunction := impl.runtimeFunction.(type) {
	case runtime_t.Function:
		for index, value := range impl.refs {
			if value.GetType() == script.ValueTypeInterface && value.GetInterface() == nil {
				panic(fmt.Sprintf("Can not find ref value for '%v'", runtimeFunction.GetRefVars()[index]))
			}
		}
	default:
		// native function can't invoke as new operator
	}

	return impl.scriptContext.(runtime.ScriptInterpreter).InvokeNew(impl.GetOwner(), args...)
}

func (impl *Component) Init() {
	if runtimeFunction, ok := impl.runtimeFunction.(runtime_t.Function); ok {
		scriptContext := impl.scriptContext.(runtime.ScriptContext)

		impl.refNames = runtimeFunction.GetRefVars()
		refLength := len(impl.refNames)
		impl.refs = make([]*script.Value, refLength)

		for i := 0; i < refLength; i++ {
			scriptContext.GetRefByName(impl.refNames[i], &impl.refs[i])
		}

		impl.memberCaches = make([]*list.List, len(runtimeFunction.GetMembers()))
		for i := range impl.memberCaches {
			impl.memberCaches[i] = list.New()
		}
	}
}

func (impl *Component) getMemberNames() []string {
	return impl.runtimeFunction.(runtime_t.Function).GetMembers()
}

func NewScriptFunction(owner, runtimeFunction, ctx interface{}) *Component {
	ret := &Component{
		ComponentType:   script.MakeComponentType(owner),
		runtimeFunction: runtimeFunction,
		scriptContext:   ctx,
		Component:       object.NewScriptObject(owner, ctx, 0),
	}

	switch f := runtimeFunction.(type) {
	case runtime_t.Function:
		ret.scriptFunction = true
		ret.scriptRuntimeFunction = f
		f.RegisterFunction(uintptr(unsafe.Pointer(ret)))
		ret._compRef = &ComopnentRef{
			pointer: ^uintptr(unsafe.Pointer(ret)),
			f:       f,
		}

		rt.SetFinalizer(ret._compRef, (*ComopnentRef).Dispose)
	case runtime_t.NativeFunction:
		ret.nativeRuntimeFunction = f
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
