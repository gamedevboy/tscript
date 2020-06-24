package _break

import (
    "container/list"
    "strings"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
    "tklibs/script/opcode"
)

type Component struct {
    debug.Component
    script.ComponentType
}

func (c *Component) Format(ident int, formatBuilder *strings.Builder) {
    panic("implement me")
}

var _ ast.Statement = &Component{}

func NewBreak(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}

func (*Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)
    ret := _func.GetInstructionList().Back()

    _func.GetBreakList().PushBack(_func.AddInstructionABx(opcode.JumpTo, opcode.Flow, compiler.NewSmallIntOperand(-1),
        compiler.NewIntOperand(0)).Value)

    return ret
}
