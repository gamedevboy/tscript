package expression

import "tklibs/script/compiler/ast"

type Member interface {
    ast.Expression
    GetLeft() interface{}
    GetRight() interface{}
}
