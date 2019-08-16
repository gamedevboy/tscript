package runtime

import "tklibs/script"

type Scope interface {
    GetFunction() interface{}
    GetArgList() []script.Value
    GetLocalVarList() []script.Value
    AddToRefList(value *script.Value, valuePtr **script.Value)
    KeepRefs()
}
