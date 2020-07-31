package object

import (
	"strings"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/opcode"
)

type Component struct {
	script.ComponentType
	values []expression.ObjectEntry
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	panic("implement me")
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

	for _, value := range impl.values {
		index := _func.GetIndexOfMemberList(value.Name)
		if index == -1 {
			index = _func.GetMemberList().Len()
			_func.GetMemberList().PushBack(value.Name)
		}
		_func.AddInstructionABC(opcode.StoreField, opcode.Memory, r, compiler.NewSmallIntOperand(int16(index)),
			value.Function.(ast.Expression).Compile(f, nil))
	}

	return r
}

func (impl *Component) GetKeyValueMap() *[]expression.ObjectEntry {
	return &impl.values
}

var _ expression.Object = &Component{}

func NewObject(owner interface{}) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
		values:        make([]expression.ObjectEntry, 0),
	}
}
