package while

import (
    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
)

type WhileStatementComponent struct {
    *whileImplement
}

func NewWhileStatementComponent(owner interface{}) *WhileStatementComponent {
    return &WhileStatementComponent{
        &whileImplement{
            ComponentType: script.NewComponentType(owner),
        },
    }
}

type whileImplement struct {
    *script.ComponentType
    condition interface{}
    body      interface{}
}

func (impl *whileImplement) SetCondition(value interface{}) {
    impl.condition = value
}

func (impl *whileImplement) SetBody(value interface{}) {
    impl.body = value
}

func (impl *whileImplement) Compile(f interface{}) {
    _func := f.(compiler.Function)

    _func.PushBreakList()

    start := _func.AddInstruction(script.OpCodeNop, nil)

    impl.condition.(ast.Node).Compile(f)

    jmp := _func.AddInstruction(script.OpCodeJEZ, nil)

    impl.body.(ast.Node).Compile(f)
    _func.AddInstruction(script.OpCodeJmp, uint32(start.Value.(*ast.Instruction).Index))

    end := _func.AddInstruction(script.OpCodeNop, nil)

    breakPos := uint32(end.Value.(*ast.Instruction).Index)

    jmp.Value.(*ast.Instruction).Operand = breakPos
    for it := _func.GetBreakList().Front(); it != nil; it = it.Next() {
        it.Value.(*ast.Instruction).Operand = breakPos
    }
    _func.PopBreakList()
}

func (impl *whileImplement) GetCondition() interface{} {
    return impl.condition
}

func (impl *whileImplement) GetBody() interface{} {
    return impl.body
}
