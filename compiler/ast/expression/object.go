package expression

import "tklibs/script/compiler/ast"

type ObjectEntry struct {
	Name     string
	Function interface{}
}

type Object interface {
	ast.Expression
	GetKeyValueMap() *[]ObjectEntry
}
