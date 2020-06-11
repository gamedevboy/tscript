package statement

import (
    "container/list"

    "tklibs/script/compiler/ast"
)

type Switch interface {
    ast.Statement
    GetTargetValue() interface{}
    SetTargetValue(value interface{})

    GetCaseList() *list.List

    GetDefaultCase() interface{}
    SetDefaultCase(value interface{})
}
