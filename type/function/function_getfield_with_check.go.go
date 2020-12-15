// +build check_cross_write

package function

import (
    "tklibs/script"
    "tklibs/script/runtime"
)

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
        runtimeObj, ok := obj.(runtime.Object)

        if ok && runtimeObj.GetRuntimeTypeInfo().(runtime.TypeInfo).GetContext() != impl.
            scriptContext {
            panic("cross context get field")
        }

        offset := impl.getFieldCache(obj, index).offset
        if offset > -1 {
            return *obj.(runtime.Object).GetByIndex(offset)
        }

        return target.ScriptGet(impl.getMemberNames()[index])
    }

    return script.NullValue
}
