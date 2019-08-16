package _map

import (
    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/opcode"
)

type Component struct {
    script.ComponentType
    values map[string]interface{}
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
    _func := f.(compiler.Function)
    if r == nil {
        r = compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABx(opcode.Map, opcode.Memory, r, compiler.NewIntOperand(script.Int(len(impl.values))))
    } else {
        n := compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABx(opcode.Map, opcode.Memory, n, compiler.NewIntOperand(script.Int(len(impl.values))))
        _func.AddInstructionABx(opcode.Move, opcode.Memory, r, n)
    }

    for varName, v := range impl.values {
        index := _func.GetIndexOfMemberList(varName)
        if index == -1 {
            index = _func.GetMemberList().Len()
            _func.GetMemberList().PushBack(varName)
        }
        _func.AddInstructionABC(opcode.StoreField, opcode.Memory, r, compiler.NewSmallIntOperand(int16(index)),
            v.(ast.Expression).Compile(f, nil))
    }

    return r
}

func (impl *Component) GetKeyValueMap() map[string]interface{} {
    return impl.values
}

var _ expression.Map = &Component{}

func NewMap(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        values:        make(map[string]interface{}),
    }
}
