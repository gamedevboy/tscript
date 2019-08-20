package while

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    _const "tklibs/script/compiler/ast/expression/const"
    "tklibs/script/compiler/debug"
    "tklibs/script/opcode"
)

type Component struct {
    debug.Component
    script.ComponentType
    condition interface{}
    body      interface{}
}

func NewWhileStatementComponent(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}

func (impl *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    _func.PushBreakList()

    start := _func.GetInstructionList().Back()

    var r *compiler.Operand

    if impl.condition != nil {
        r = impl.condition.(ast.Expression).Compile(f, nil)
    } else {
        tc := &struct {
            *_const.Component
        }{}
        tc.Component = _const.NewConst(tc, true)
        r = tc.Compile(f, nil)
    }

    if start != nil {
        start = start.Next()
    } else {
        start = _func.GetInstructionList().Front()
    }

    jmp := _func.AddInstructionABx(opcode.JumpWhenFalse, opcode.Flow, r, compiler.NewIntOperand(0))

    impl.body.(ast.Statement).Compile(f)

    _func.AddInstructionABx(opcode.Jump, opcode.Flow, compiler.NewSmallIntOperand(-1),
        compiler.NewIntOperand(script.Int(start.Value.(*ast.Instruction).Index)))

    end := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

    breakPos := script.Int(end.Value.(*ast.Instruction).Index + 1)

    jmp.Value.(*ast.Instruction).GetABx().B = breakPos

    for it := _func.GetBreakList().Front(); it != nil; it = it.Next() {
        it.Value.(*ast.Instruction).GetABx().B = breakPos
    }

    _func.PopBreakList()

    return start
}

func (impl *Component) GetCondition() interface{} {
    return impl.condition
}

func (impl *Component) SetCondition(value interface{}) {
    impl.condition = value
}

func (impl *Component) GetBody() interface{} {
    return impl.body
}

func (impl *Component) SetBody(value interface{}) {
    impl.body = value
}
