package expression

import "tklibs/script/compiler/ast"

type Array interface {
    ast.Expression
    GetArgListExpression() interface{}
}
