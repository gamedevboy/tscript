package _if

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/debug"
    "tklibs/script/compiler/token"
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

func (impl *Component) expandConditionExpression(e interface{}, conditionList *list.List) {
    switch val := e.(type) {
    case expression.Binary:
        switch val.GetOpType() {
        case token.TokenTypeLAND:
            if conditionList.Len() == 0 {
                conditionList.PushBack(true)
            }
            impl.expandConditionExpression(val.GetLeft(), conditionList)
            conditionList.PushBack(true)
            conditionList.PushBack(val.GetRight())
        case token.TokenTypeLOR:
            if conditionList.Len() == 0 {
                conditionList.PushBack(false)
            }
            impl.expandConditionExpression(val.GetLeft(), conditionList)
            conditionList.PushBack(false)
            conditionList.PushBack(val.GetRight())
        default:
            conditionList.PushBack(val)
        }
    default:
        conditionList.PushBack(val)
    }
}

func (impl *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    ret := _func.GetInstructionList().Back()

    jumpList := list.New() // jump to the end
    skipJumpList := list.New() // jump to the body start

    conditionList := list.New()

    impl.expandConditionExpression(impl.condition, conditionList)

    r := compiler.NewRegisterOperand(_func.AllocRegister(""))
    lastOp := false
    for it := conditionList.Front(); it != nil; it = it.Next() {
        switch val := it.Value.(type) {
        case bool:
            lastOp = val
        case ast.Expression:
            val.Compile(f, r)

            if lastOp {
                jumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, r,
                    compiler.NewIntOperand(0)))
            } else {
                skipJumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenTrue, opcode.Flow, r,
                    compiler.NewIntOperand(0)))
            }
        }


        // if it.Value.(bool) {
        //     it.Next().Value.(ast.Expression).Compile(f, r)
        //
        //     next := it.Next().Next()
        //     if next != nil && next.Value.(bool) {
        //         jumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, r,
        //             compiler.NewIntOperand(0)))
        //     }
        // } else {
        //     skipJumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenTrue, opcode.Flow, r,
        //         compiler.NewIntOperand(0)))
        //     next := it.Next().Value.(ast.Expression).Compile(f, nil)
        //     _func.AddInstructionABC(opcode.LogicOr, opcode.Logic, r, r, next)
        //     jumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, r,
        //         compiler.NewIntOperand(0)))
        // }
    }

    if !lastOp {
        jumpList.PushBack(_func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, r,
            compiler.NewIntOperand(0)))
    }

    bodyStart := impl.body.(ast.Statement).Compile(f)

    for it := skipJumpList.Front(); it != nil; it= it.Next() {
        it.Value.(*list.Element).Value.(*ast.Instruction).GetABx().B = bodyStart.Value.(*ast.Instruction).Index
    }

    endJmp := _func.AddInstructionABx(opcode.Jump, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    jumpTarget := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    if impl.elseBody != nil {
        impl.elseBody.(ast.Statement).Compile(f)
    }

    endTarget := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    for it := jumpList.Front(); it != nil; it = it.Next() {
        it.Value.(*list.Element).Value.(*ast.Instruction).GetABx().B = jumpTarget.Value.(*ast.Instruction).Index + 1
    }

    endJmp.Value.(*ast.Instruction).GetABx().B = endTarget.Value.(*ast.Instruction).Index + 1

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
