// +build check_cross_write

package function

import (
"tklibs/script"
"tklibs/script/runtime"
)

const (
    CrossWriteCheck = true
)

func (impl *Component) SetFieldByMemberIndex(obj interface{}, index script.Int, value script.Value) {
    runtimeObj, ok := obj.(runtime.Object)

    if ok && runtimeObj.GetRuntimeTypeInfo().(runtime.TypeInfo).GetContext() != impl.
        scriptContext {
        panic("cross context set field")
    }

    switch target := obj.(type) {
    case script.Map:
        target.Set(script.String(impl.getMemberNames()[index]), value.Get())
    default:
        offset := impl.getFieldCache(runtimeObj, index).offset
        if offset > -1 {
            runtimeObj.SetByIndex(offset, value)
        } else {
            runtimeObj.(script.Object).ScriptSet(impl.getMemberNames()[index], value)
        }
    }
}