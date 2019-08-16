package array

import (
    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/opcode"
)

type ArrayExpressionComponent struct {
    *arrayImplement
}

type arrayImplement struct {
    script.ComponentType
    argList interface{}
}

func (impl *arrayImplement) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
    _func := f.(compiler.Function)
    argList := impl.argList.(expression.ArgList)

    if r == nil {
        r = compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABx(opcode.Array, opcode.Memory, r, compiler.NewIntOperand(script.Int(argList.GetExpressionList().Len())))
    } else {
        n := compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABx(opcode.Array, opcode.Memory, n, compiler.NewIntOperand(script.Int(argList.GetExpressionList().Len())))
        _func.AddInstructionABx(opcode.Move, opcode.Memory, r, n)
    }

    i := int16(0)
    for it := argList.GetExpressionList().Front(); it != nil; it = it.Next() {
        _func.AddInstructionABC(opcode.StoreElement, opcode.Memory, r, compiler.NewSmallIntOperand(i),
            it.Value.(ast.Expression).Compile(f, nil))
        i++
    }

    return r
}

func (impl *arrayImplement) GetArgListExpression() interface{} {
    return impl.argList
}

func NewArrayExpression(owner, arglist interface{}) *ArrayExpressionComponent {
    return &ArrayExpressionComponent{
        &arrayImplement{
            ComponentType: script.MakeComponentType(owner),
            argList:       arglist,
        },
    }
}
