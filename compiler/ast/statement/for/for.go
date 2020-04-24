package _for

import (
	"container/list"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/compiler/ast/expression/const"
	"tklibs/script/compiler/debug"
	"tklibs/script/compiler/token"
	"tklibs/script/opcode"
)

type Component struct {
	debug.Component
	script.ComponentType
	init      interface{}
	condition interface{}
	step      interface{}
	body      interface{}
}

func NewForStatementComponent(owner interface{}) *Component {
	return &Component{ComponentType: script.MakeComponentType(owner)}
}

func (impl *Component) expandConditionExpression(e interface{}, conditionList *list.List) {
	switch val := e.(type) {
	case expression.Binary:
		switch val.GetOpType() {
		case token.TokenTypeLAND:
			impl.expandConditionExpression(val.GetLeft(), conditionList)
			conditionList.PushBack(true)
			impl.expandConditionExpression(val.GetRight(), conditionList)
		case token.TokenTypeLOR:
			impl.expandConditionExpression(val.GetLeft(), conditionList)
			conditionList.PushBack(false)
			impl.expandConditionExpression(val.GetRight(), conditionList)
		default:
			conditionList.PushBack(e)
		}
	default:
		conditionList.PushBack(e)
	}
}

func (impl *Component) Compile(f interface{}) *list.Element {
	_func := f.(compiler.Function)

	_func.PushBreakList()
	defer _func.PopBreakList()

	_func.PushContinueList()
	defer _func.PopContinueList()

	if impl.init != nil {
		impl.init.(ast.Statement).Compile(f)
	}

	start := _func.GetInstructionList().Back()

	r := compiler.NewRegisterOperand(_func.AllocRegister(""))

	jumpList := list.New()
	skipJumpList := list.New()

	if impl.condition != nil {
		conditionList := list.New()
		conditionList.PushBack(true)

		impl.expandConditionExpression(impl.condition, conditionList)

		for it := conditionList.Front(); it != nil; {
			if it.Value.(bool) {
				it.Next().Value.(ast.Expression).Compile(f, r)

				next := it.Next().Next()
				if next != nil && next.Value.(bool) {
					jumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r, compiler.NewIntOperand(0)))
				}
			} else {
				skipJumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r, compiler.NewIntOperand(0)))
				next := it.Next().Value.(ast.Expression).Compile(f, nil)
				_func.AddInstructionABC(opcode.LogicOr, opcode.Logic, r, r, next)
				jumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r, compiler.NewIntOperand(0)))
			}

			it = it.Next().Next()
		}
	} else {
		tc := &struct {
			*_const.Component
		}{}
		tc.Component = _const.NewConst(tc, true)
		tc.Compile(f, r)
	}

	if start != nil {
		start = start.Next()
	} else {
		start = _func.GetInstructionList().Front()
	}

	jmp := _func.AddInstructionABx(opcode.Jump, opcode.Flow, r, compiler.NewIntOperand(0))

	bodyStart := impl.body.(ast.Statement).Compile(f)

	for it := skipJumpList.Front(); it != nil; it = it.Next() {
		it.Value.(*list.Element).Value.(*ast.Instruction).GetABx().B = -bodyStart.Value.(*ast.Instruction).Index
	}

	continuePos := -1

	if impl.step != nil {
		continueInst := impl.step.(ast.Statement).Compile(f)
		continuePos = int(continueInst.Value.(*ast.Instruction).Index)
	}

	startPos := start.Value.(*ast.Instruction).Index

	_func.AddInstructionABx(opcode.JumpTo, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(script.Int(startPos)))

	end := _func.AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

	breakPos := end.Value.(*ast.Instruction).Index + 1

	for it := jumpList.Front(); it != nil; it = it.Next() {
		it.Value.(*list.Element).Value.(*ast.Instruction).GetABx().B = breakPos
	}

	jmp.Value.(*ast.Instruction).GetABx().B = breakPos

	for it := _func.GetBreakList().Front(); it != nil; it = it.Next() {
		it.Value.(*ast.Instruction).GetABx().B = breakPos
	}

	for it := _func.GetContinueList().Front(); it != nil; it = it.Next() {
		if continuePos == -1 {
			it.Value.(*ast.Instruction).GetABx().B = startPos
		} else {
			it.Value.(*ast.Instruction).GetABx().B = int32(continuePos)
		}
	}

	return start
}

func (impl *Component) GetInit() interface{} {
	return impl.init
}

func (impl *Component) SetInit(value interface{}) {
	impl.init = value
}

func (impl *Component) GetCondition() interface{} {
	return impl.condition
}

func (impl *Component) SetCondition(value interface{}) {
	impl.condition = value
}

func (impl *Component) GetStep() interface{} {
	return impl.step
}

func (impl *Component) SetStep(value interface{}) {
	impl.step = value
}

func (impl *Component) GetBody() interface{} {
	return impl.body
}

func (impl *Component) SetBody(value interface{}) {
	impl.body = value
}
