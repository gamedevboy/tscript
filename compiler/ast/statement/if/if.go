package _if

import (
    "container/list"

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
    condition interface{}
    body      interface{}
    elseBody  interface{}
}

var _ statement.If = &Component{}

func (impl *Component) SetCondition(value interface{}) {
    impl.condition = value
}

func (impl *Component) SetBody(value interface{}) {
    impl.body = value
}

func (impl *Component) SetElseBody(value interface{}) {
    impl.elseBody = value
}

func (impl *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    ret := _func.GetInstructionList().Back()

    jump := _func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, impl.condition.(ast.Expression).Compile(f, nil),
        compiler.NewIntOperand(0))

    impl.body.(ast.Statement).Compile(f)
    endJmp := _func.AddInstructionABx(opcode.Jump, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    jumpTarget := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    if impl.elseBody != nil {
        impl.elseBody.(ast.Statement).Compile(f)
    }

    endTarget := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    jump.Value.(*ast.Instruction).GetABx().B = script.Int(jumpTarget.Value.(*ast.Instruction).Index + 1)
    endJmp.Value.(*ast.Instruction).GetABx().B = script.Int(endTarget.Value.(*ast.Instruction).Index + 1)

    return ret
}

func (impl *Component) GetCondition() interface{} {
    return impl.condition
}

func (impl *Component) GetBody() interface{} {
    return impl.body
}

func (impl *Component) GetElseBody() interface{} {
    return impl.elseBody
}

func NewIfStatementComponent(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}
