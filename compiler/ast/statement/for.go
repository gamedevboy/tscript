package statement

import (
    "tklibs/script/compiler/ast"
)

type For interface {
    ast.Statement
    GetInit() interface{}
    SetInit(interface{})
    GetCondition() interface{}
    SetCondition(interface{})
    GetStep() interface{}
    SetStep(interface{})
    GetBody() interface{}
    SetBody(interface{})
}
