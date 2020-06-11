package statement

import (
	"container/list"

	"tklibs/script/compiler/ast"
)

type Block interface {
	ast.Statement
	GetStatementList() *list.List
}

