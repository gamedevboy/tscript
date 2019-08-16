package expression

import "tklibs/script/compiler/ast"

type Function interface {
    ast.Expression
    GetArgList() interface{}
    SetArgList(interface{})
    GetBlock() interface{}
    SetBlock(interface{})
    SetMetaIndex(uint32)
    GetMetaIndex() uint32
    GetName() string
    SetName(string)
    SetCaptureThis(val bool)
    GetCaptureThis() bool
}
