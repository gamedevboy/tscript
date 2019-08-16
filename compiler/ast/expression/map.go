package expression

import "tklibs/script/compiler/ast"

type Map interface {
    ast.Expression
    GetKeyValueMap() map[string]interface{}
}
