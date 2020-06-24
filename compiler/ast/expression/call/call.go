package call

import (
	"container/list"
	"fmt"
	"strings"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/compiler/token"
	"tklibs/script/opcode"
)

type Component struct {
	script.ComponentType
	expression interface{}
	argList    interface{}
	isNew      bool
	option     bool
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	if impl.isNew {
		formatBuilder.WriteString("new ")
	}
	impl.GetExpression().(ast.Node).Format(ident, formatBuilder)
	formatBuilder.WriteString("(")
	impl.GetArgList().(ast.Node).Format(ident, formatBuilder)
	formatBuilder.WriteString(")")
}

var _ expression.Call = &Component{}

func (impl *Component) SetNew(value bool) {
	impl.isNew = value
}

func (impl *Component) GetExpression() interface{} {
	return impl.expression
}

func (impl *Component) GetNew() bool {
	return impl.isNew
}

func (impl *Component) String() string {
	return fmt.Sprint(impl.expression, "(", impl.argList, ")")
}

func (impl *Component) GetArgList() interface{} {
	return impl.argList
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
	argList := impl.argList.(expression.ArgList).GetExpressionList()

	_func := f.(compiler.Function)

	rf := impl.expression.(ast.Expression).Compile(f, nil)

	if impl.option {
		rn := _func.AllocRegister("")
		_func.AddInstructionABx(opcode.LoadNil, opcode.Const, compiler.NewRegisterOperand(rn), nil)
		_func.AddInstructionABC(opcode.Equal, opcode.Logic, compiler.NewRegisterOperand(rn), rf, compiler.NewRegisterOperand(rn))
		jump := _func.AddInstructionABx(opcode.Jump, opcode.Flow, compiler.NewRegisterOperand(rn), compiler.NewIntOperand(0))
		end := call(f, r, _func, impl, argList, rf)
		endJump := _func.AddInstructionABx(opcode.JumpTo, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))
		nilEnd := _func.AddInstructionABx(opcode.LoadNil, opcode.Const, end, nil)
		nop := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

		jump.Value.(*ast.Instruction).GetABx().B = -nilEnd.Value.(*ast.Instruction).Index
		endJump.Value.(*ast.Instruction).GetABx().B = nop.Value.(*ast.Instruction).Index

		return end
	} else {
		return call(f, r, _func, impl, argList, rf)
	}
}

func call(f interface{}, r *compiler.Operand, _func compiler.Function, impl *Component, argList *list.List, rf *compiler.Operand) *compiler.Operand {
	regCount := _func.GetRegisterCount()

	_ = _func.AllocRegister("")                                  // reverse for return value
	this := compiler.NewRegisterOperand(_func.AllocRegister("")) // reverse for this value

	_func.PushRegisters()

	switch me := impl.expression.(type) {
	case expression.Member:
		if me.GetLeft() != nil {
			me.GetLeft().(ast.Expression).Compile(f, this)
		} else {
			_func.AddInstructionABx(opcode.LoadNil, opcode.Const, this, nil)
		}
	}

	_func.ReleaseAllRegisters()

	argLen := int16(argList.Len())

	for it := argList.Front(); it != nil; it = it.Next() {
		r := compiler.NewRegisterOperand(_func.AllocRegister(""))
		_func.PushRegisters()

		if u, ok := it.Value.(expression.Unary); ok && u.GetTokenType() == token.TokenTypeELLIPSIS && it.Next() == nil {
			argLen = -argLen
			u.GetExpression().(ast.Expression).Compile(f, r)
		} else {
			it.Value.(ast.Expression).Compile(f, r)
		}

		_func.PopRegisters()
	}

	if impl.isNew {
		_func.AddInstructionABC(opcode.Call, opcode.Flow, rf, compiler.NewSmallIntOperand(regCount<<1+1),
			compiler.NewSmallIntOperand(argLen))
	} else {
		_func.AddInstructionABC(opcode.Call, opcode.Flow, rf, compiler.NewSmallIntOperand(regCount<<1),
			compiler.NewSmallIntOperand(argLen))
	}

	_func.PopRegisters()

	if r == nil {
		return compiler.NewRegisterOperand(&compiler.Register{Index: regCount})
	} else {
		_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRegisterOperand(&compiler.Register{Index: regCount}))
		return r
	}
}

func NewCall(owner, e, argList interface{}, option bool) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
		expression:    e,
		argList:       argList,
		option:        option,
	}
}
