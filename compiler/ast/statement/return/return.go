package _return

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
    "tklibs/script/opcode"
)

type Component struct {
    debug.Component
    script.ComponentType
    expression interface{}
}

func NewReturn(owner interface{}) *Component {
    return &Component{ComponentType: script.MakeComponentType(owner)}
}

func (impl *Component) Compile(f interface{}) *list.Element {
    if impl.expression != nil {
        impl.expression.(ast.Expression).Compile(f, compiler.NewRegisterOperand(&compiler.Register{Index: 0}))
    }

    return f.(compiler.Function).AddInstructionABx(opcode.Ret, opcode.Flow, compiler.NewRegisterOperand(&compiler.Register{Index: 0}), compiler.NewIntOperand(0))
}

func (impl *Component) GetExpression() interface{} {
    return impl.expression
}

func (impl *Component) SetExpression(value interface{}) {
    impl.expression = value
}