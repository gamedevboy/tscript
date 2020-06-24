package expression

import (
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/token"
)

type Binary interface {
    ast.Expression
    GetLeft() interface{}
    GetRight() interface{}
    GetOpType() token.TokenType
    SetParen(value bool)
    GetParen() bool
}
