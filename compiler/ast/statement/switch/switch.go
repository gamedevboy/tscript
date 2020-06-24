package _switch

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

func NewSwitch(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}

type Component struct {
    debug.Component
    script.ComponentType
    targetValue interface{}
    caseList    list.List
    defaultCase interface{}
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
    panic("implement me")
}

var _ ast.Statement = &Component{}

func (impl *Component) GetTargetValue() interface{} {
    return impl.targetValue
}

func (impl *Component) SetTargetValue(value interface{}) {
    impl.targetValue = value
}

func (impl *Component) GetCaseList() *list.List {
    return &impl.caseList
}

func (impl *Component) GetDefaultCase() interface{} {
    return impl.defaultCase
}

func (impl *Component) SetDefaultCase(value interface{}) {
    impl.defaultCase = value
}

func (impl *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    ret := _func.GetInstructionList().Back()

    _func.PushBreakList()
    defer _func.PopBreakList()

    jmpList := make([]*list.Element, impl.caseList.Len()+1)
    caseStartList := make([]*list.Element, impl.caseList.Len()+1)

    target := impl.targetValue.(ast.Expression).Compile(f, nil)
    _func.PushRegisters()
    defer _func.PopRegisters()

    // compile case condition
    index := 0
    for it := impl.caseList.Front(); it != nil; it = it.Next() {
        v := it.Value.(statement.Case).GetValue().(ast.Expression).Compile(f, nil)
        r := compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABC(opcode.NotEqual, opcode.Logic, r, target, v)
        jmpList[index] = _func.AddInstructionABx(opcode.Jump, opcode.Flow, r, compiler.NewIntOperand(0))
        index++
    }

    jmpList[index] = _func.AddInstructionABx(opcode.JumpTo, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    index = 0
    for it := impl.caseList.Front(); it != nil; it = it.Next() {
        caseStartList[index] = _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))
        it.Value.(statement.Case).GetBlock().(ast.Statement).Compile(f)
        index++
    }

    // compile default case
    caseStartList[index] = _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))
    if impl.defaultCase != nil {
        impl.defaultCase.(ast.Statement).Compile(f)
    }

    for idx, jmp := range jmpList {
        jmp.Value.(*ast.Instruction).GetABx().B = caseStartList[idx].Value.(*ast.Instruction).Index
    }

    end := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))
    breakPos := end.Value.(*ast.Instruction).Index
    for it := _func.GetBreakList().Front(); it != nil; it = it.Next() {
        it.Value.(*ast.Instruction).GetABx().B = breakPos
    }

    return ret
}
