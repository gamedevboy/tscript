package object

import (
    "sort"

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
        _func.AddInstructionABx(opcode.Object, opcode.Memory, r, compiler.NewIntOperand(script.Int(len(impl.values))))
    } else {
        n := compiler.NewRegisterOperand(_func.AllocRegister(""))
        _func.AddInstructionABx(opcode.Object, opcode.Memory, n, compiler.NewIntOperand(script.Int(len(impl.values))))
        _func.AddInstructionABx(opcode.Move, opcode.Memory, r, n)
    }

    keys := make([]string, len(impl.values))

    idx := 0
    for varName := range impl.values {
        keys[idx] = varName
        idx++
    }

    sort.Strings(keys)

    for _, varName := range keys {
        v := impl.values[varName]
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

var _ expression.Object = &Component{}

func NewObject(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        values:        make(map[string]interface{}),
    }
}
