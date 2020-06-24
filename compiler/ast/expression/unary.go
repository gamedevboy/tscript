package expression

import (
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/token"
)

type Unary interface {
    ast.Expression
    GetTokenType() token.TokenType
    GetExpression() interface{}
}
