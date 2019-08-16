package expression

import "tklibs/script/compiler/ast"

type Const interface {
    ast.Expression
    GetValue() interface{}
}
