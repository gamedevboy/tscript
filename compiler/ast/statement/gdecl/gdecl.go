package decl

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
)

type Component struct {
    debug.Component
    *script.ComponentType
    name       string
    expression interface{}
}

func (ds *Component) GetName() string {
    return ds.name
}

func (ds *Component) SetName(name string) {
    ds.name = name
}

func (ds *Component) GetExpression() interface{} {
    return ds.expression
}

func (ds *Component) SetExpression(expression interface{}) {
    ds.expression = expression
}

func (ds *Component) String() string {
    if ds.expression != nil {
        return fmt.Sprint("gvar ", ds.name, "=", ds.expression)
    }

    return fmt.Sprint("gvar ", ds.name)
}

func (ds *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    ret := _func.GetInstructionList().Back()

    if ds.expression != nil {
        index := _func.GetIndexOfLocalList(ds.name)
        ds.expression.(ast.Expression).Compile(f, compiler.NewRegisterOperand(_func.GetRegisterByLocalIndex(index)))
        if ret == nil {
            return nil
        }
        return ret.Next()
    }

    return ret
}

func NewGlobalDecl(owner interface{}) *Component {
    return &Component{Component: script.NewComponentType(owner)}
}
