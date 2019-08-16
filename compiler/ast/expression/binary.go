package expression

import "tklibs/script/compiler/ast"

type Binary interface {
    ast.Expression
    GetLeft() interface{}
    GetRight() interface{}
}
