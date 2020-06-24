package expression

import (
    "container/list"
    "fmt"
    "strings"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
    "tklibs/script/opcode"
)

var _ ast.Statement = &Component{}

type Component struct {
    debug.Component
    script.ComponentType
    expression interface{}
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
    impl.expression.(ast.Node).Format(ident, formatBuilder)
}

func (impl *Component) String() string {
    return fmt.Sprint(impl.expression)
}

func (impl *Component) GetExpression() interface{} {
    return impl.expression
}

func (impl *Component) SetExpression(expression interface{}) {
    impl.expression = expression
}

func (impl *Component) Compile(f interface{}) *list.Element {
    cur := f.(compiler.Function).GetInstructionList().Back()

    if impl.expression != nil {
        impl.expression.(ast.Expression).Compile(f, nil)
    } else {
        return f.(compiler.Function).AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1),
            compiler.NewIntOperand(0))
    }

    if cur == nil {
        return f.(compiler.Function).GetInstructionList().Front()
    }

    return cur.Next()
}

func NewExpressionStatement(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}
