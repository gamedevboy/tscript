package statement

import (
	"tklibs/script/compiler/ast"
)

type Expression interface {
	ast.Statement
	GetExpression() interface{}
	SetExpression(interface{})
}

