package function

import (
    "container/list"
    "fmt"
    "sync/atomic"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/type/object"
)

type Component struct {
    script.ComponentType
    *object.Component
    scriptContext   interface{}
    runtimeFunction interface{}
    refs            []*script.Value
    memberNames     []string
    memberCaches    []*list.List
    initialized     int32
    this            script.Value
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
        if int(index) < len(impl.memberCaches) {
            for it := impl.memberCaches[index].Front(); it != nil; it = it.Next() {
                if it.Value.(*fieldCache).cacheType == runtimeObj.GetRuntimeTypeInfo() {
                    return it.Value.(*fieldCache)
                }
            }
        } else {
            impl.memberCaches = make([]*list.List, len(impl.getMemberNames()))

            for i := range impl.memberCaches {
                impl.memberCaches[i] = list.New()
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
        key := script.Value{}
        key.Set(script.String(impl.getMemberNames()[index]))

        if target.ContainsKey(key) {
            return target.GetValue(key)
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
        k := script.Value{}
        k.Set(script.String(impl.getMemberNames()[index]))
        target.SetValue(k, value)
    default:
        offset := impl.getFieldCache(obj, index).offset
        if offset > - 1 {
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
    return impl.Init().scriptContext.(runtime.ScriptInterpreter).InvokeFunction(impl.GetOwner(), this, args...)
}

func (impl *Component) New(args ...interface{}) interface{} {
    impl.Init()
    switch runtimeFunction := impl.runtimeFunction.(type) {
    case runtime.Function:
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

func (impl *Component) Init() *Component {
    if atomic.CompareAndSwapInt32(&impl.initialized, 0, 1) {
        runtimeFunction := impl.runtimeFunction.(runtime.Function)
        scriptContext := impl.scriptContext.(runtime.ScriptContext)

        for i := 0; i < len(runtimeFunction.GetRefVars()); i++ {
            scriptContext.GetRefByName(runtimeFunction.GetRefVars()[i], &impl.refs[i])
        }
    }

    return impl
}

func (impl *Component) getMemberNames() []string {
    return impl.runtimeFunction.(runtime.Function).GetMembers()
}

func NewScriptFunction(owner, runtimeFunction, ctx interface{}) *Component {
    ret := &Component{
        ComponentType:   script.MakeComponentType(owner),
        runtimeFunction: runtimeFunction,
        scriptContext:   ctx,
        Component:       object.NewScriptObject(owner, ctx, 0),
    }
    ret.SetPrototype(script.InterfaceToValue(ctx.(runtime.ScriptContext).GetFunctionPrototype()))

    switch f := runtimeFunction.(type) {
    case runtime.Function:
        ret.refs = make([]*script.Value, len(f.GetRefVars()))

        ret.memberCaches = make([]*list.List, len(ret.memberNames))

        for i := range ret.memberCaches {
            ret.memberCaches[i] = list.New()
        }
    }

    return ret
}
