package expression

import (
    "container/list"

    "tklibs/script/compiler/ast"
)

type ArgList interface {
    ast.Expression
    GetExpressionList() *list.List
}
