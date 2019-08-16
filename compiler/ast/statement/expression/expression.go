package expression

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
)

var _ ast.Statement = &Component{}

type Component struct {
    debug.Component
    script.ComponentType
    expression interface{}
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
    impl.expression.(ast.Expression).Compile(f, nil)

    if cur == nil {
        return cur
    }

    return cur.Next()
}

func NewExpressionStatement(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}
