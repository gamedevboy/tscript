package statement

import (
    "tklibs/script/compiler/ast"
)

type Case interface {
    ast.Statement
    GetValue() interface{}
    SetValue(value interface{})

    GetBlock() interface{}
    SetBlock(value interface{})
}
