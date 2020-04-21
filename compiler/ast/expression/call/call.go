package call

import (
	"fmt"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/opcode"
)

type Component struct {
	script.ComponentType
	expression interface{}
	argList    interface{}
	isNew      bool
	option     bool
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

	for it := argList.Front(); it != nil; it = it.Next() {
		r := compiler.NewRegisterOperand(_func.AllocRegister(""))
		_func.PushRegisters()
		it.Value.(ast.Expression).Compile(f, r)
		_func.PopRegisters()
	}

	if impl.isNew {
		_func.AddInstructionABC(opcode.Call, opcode.Flow, rf, compiler.NewSmallIntOperand(regCount<<1+1),
			compiler.NewSmallIntOperand(int16(argList.Len())))
	} else {
		_func.AddInstructionABC(opcode.Call, opcode.Flow, rf, compiler.NewSmallIntOperand(regCount<<1),
			compiler.NewSmallIntOperand(int16(argList.Len())))
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
