// +build !check_cross_read

package function

import (
    "tklibs/script"
    "tklibs/script/runtime"
)

const (
    CrossReadCheck = false
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
        offset := impl.getFieldCache(obj, index).offset
        if offset > -1 {
            return *obj.(runtime.Object).GetByIndex(offset)
        }

        return target.ScriptGet(impl.getMemberNames()[index])
    }

    return script.NullValue
}
