package _return

import (
    "container/list"
    "strings"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/debug"
    "tklibs/script/opcode"
)

type Component struct {
    debug.Component
    script.ComponentType
    expression interface{}
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
    formatBuilder.WriteString("return")
    if impl.expression != nil {
        formatBuilder.WriteString(" ")
        impl.expression.(ast.Node).Format(ident, formatBuilder)
    }
}

var _ statement.Return = &Component{}

func NewReturn(owner interface{}) *Component {
    return &Component{ComponentType: script.MakeComponentType(owner)}
}

func (impl *Component) Compile(f interface{}) *list.Element {
    cur := f.(compiler.Function).GetInstructionList().Back()

    if impl.expression != nil {
         impl.expression.(ast.Expression).Compile(f, compiler.NewRegisterOperand(&compiler.Register{Index: 0}))
    }

    f.(compiler.Function).AddInstructionABx(opcode.Ret, opcode.Flow, compiler.NewRegisterOperand(&compiler.Register{Index: 0}), compiler.NewIntOperand(0))

    if cur == nil {
        return f.(compiler.Function).GetInstructionList().Front()
    }

    return cur.Next()
}

func (impl *Component) GetExpression() interface{} {
    return impl.expression
}

func (impl *Component) SetExpression(value interface{}) {
    impl.expression = value
}
