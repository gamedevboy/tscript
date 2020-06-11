package statement

import (
    "tklibs/script/compiler/ast"
)

type If interface {
    ast.Statement
    GetCondition() interface{}
    SetCondition(interface{})
    GetBody() interface{}
    SetBody(interface{})
    GetElseBody() interface{}
    SetElseBody(interface{})
}
