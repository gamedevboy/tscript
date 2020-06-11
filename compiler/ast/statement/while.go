package statement

import (
    "tklibs/script/compiler/ast"
)

type While interface {
    ast.Statement
    GetCondition() interface{}
    SetCondition(interface{})
    GetBody() interface{}
    SetBody(interface{})
}
