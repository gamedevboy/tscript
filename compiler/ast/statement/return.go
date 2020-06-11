package statement

import (
    "tklibs/script/compiler/ast"
)

type Return interface {
    ast.Statement
    GetExpression() interface{}
    SetExpression(value interface{})
}
