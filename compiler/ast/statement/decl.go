package statement

import (
    "tklibs/script/compiler/ast"
)

type Decl interface {
    ast.Statement
    GetName() string
    SetName(string)

    GetExpression() interface{}
    SetExpression(interface{})

    IsGlobal() bool
    SetGlobal(bool)
}

