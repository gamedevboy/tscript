package expression

import "tklibs/script/compiler/ast"

type Object interface {
    ast.Expression
    GetKeyValueMap() map[string]interface{}
}
