package _if

import (
	"container/list"
	"strings"

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

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	panic("implement me")
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

func (impl *Component) expandConditionExpression(e interface{}, conditionList *list.List, f compiler.Function) {
	switch val := e.(type) {
	case expression.Binary:
		switch val.GetOpType() {
		case token.TokenTypeLAND:
			if conditionList.Len() == 0 {
				conditionList.PushBack(true)
			}

			if conditionList.Back().Value.(bool) == true {
				impl.expandConditionExpression(val.GetRight(), conditionList, f)
				conditionList.PushBack(true)
				impl.expandConditionExpression(val.GetLeft(), conditionList, f)
			} else {
				conditionList.PushBack(val)
			}
		case token.TokenTypeLOR:
			if conditionList.Len() == 0 {
				conditionList.PushBack(false)
			}

			if conditionList.Back().Value.(bool) == false {
				impl.expandConditionExpression(val.GetRight(), conditionList, f)
				conditionList.PushBack(false)
				impl.expandConditionExpression(val.GetLeft(), conditionList, f)
			} else {
				conditionList.PushBack(val)
			}
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

	jumpList := list.New()     // jump to the end
	skipJumpList := list.New() // jump to the body start

	conditionList := list.New()

	impl.expandConditionExpression(impl.condition, conditionList, _func)

	if conditionList.Len() < 2 {
		conditionList.PushFront(true)
	}

	r := compiler.NewRegisterOperand(_func.AllocRegister(""))

	cop := conditionList.Front().Value.(bool)

	for conditionList.Len() > 0 {
		switch val := conditionList.Back().Value.(type) {
		case bool:
			if val {
				jumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r,
					compiler.NewIntOperand(0)))
			} else {
				skipJumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r,
					compiler.NewIntOperand(0)))
			}
		case ast.Expression:
			val.Compile(f, r)
		}

		conditionList.Remove(conditionList.Back())
	}

	if !cop {
		jumpList.PushBack(_func.AddInstructionABx(opcode.Jump, opcode.Flow, r,
			compiler.NewIntOperand(0)))
	}

	bodyStart := impl.body.(ast.Statement).Compile(f)

	for it := skipJumpList.Front(); it != nil; it = it.Next() {
		it.Value.(*list.Element).Value.(*ast.Instruction).GetABx().B = -bodyStart.Value.(*ast.Instruction).Index
	}

	endJmp := _func.AddInstructionABx(opcode.JumpTo, opcode.Flow, compiler.NewSmallIntOperand(-1), compiler.NewIntOperand(0))

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

func getCurrCondition(conditionStack *list.List) bool {
	c := false

	for x := conditionStack.Front(); x != nil; x = x.Next() {
		c = c && x.Value.(bool)
	}

	return c
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
