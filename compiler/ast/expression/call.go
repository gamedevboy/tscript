package expression

import "tklibs/script/compiler/ast"

type Call interface {
    ast.Expression

    SetNew(value bool)
    GetNew() bool
    GetExpression() interface{}
    GetArgList() interface{}
}
