package binary

import (
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
	left, right interface{}
	opType      token.TokenType
	paren       bool
}

func (impl *Component) SetParen(value bool) {
	impl.paren = value
}

func (impl *Component) GetParen() bool {
	return impl.paren
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	if impl.paren {
		formatBuilder.WriteString("(")
	}
	impl.left.(ast.Node).Format(ident, formatBuilder)
	formatBuilder.WriteString(fmt.Sprintf(" %v ", impl.opType))
	impl.right.(ast.Node).Format(ident, formatBuilder)
	if impl.paren {
		formatBuilder.WriteString(")")
	}
}

func (impl *Component) GetLeft() interface{} {
	return impl.left
}

func (impl *Component) GetRight() interface{} {
	return impl.right
}

var _ expression.Binary = &Component{}

func (impl *Component) String() string {
	return fmt.Sprint("(", impl.left, impl.opType, impl.right, ")")
}

func (impl *Component) GetOpType() token.TokenType {
	return impl.opType
}

func (impl *Component) compileStore(f, left interface{}, r *compiler.Operand) {
	_func := f.(compiler.Function)
	m := left.(expression.Member)

	switch varName := m.GetRight().(type) {
	case string:
		impl.compileStoreByString(m, _func, varName, f, r)
	case expression.Const:
		switch vn := varName.GetValue().(type) {
		case script.String:
			impl.compileStoreByString(m, _func, string(vn), f, r)
		default:
			_func.AddInstructionABC(opcode.StoreElement, opcode.Memory, m.GetLeft().(ast.Expression).Compile(f, nil), varName.Compile(f,
				nil), r)
		}
	case ast.Expression:
		_func.AddInstructionABC(opcode.StoreElement, opcode.Memory, m.GetLeft().(ast.Expression).Compile(f, nil), varName.Compile(f,
			nil), r)
	}
}

func (impl *Component) compileStoreByString(m expression.Member, _func compiler.Function, varName string, f interface{}, r *compiler.Operand) {
	if m.GetLeft() != nil {
		index := _func.GetIndexOfMemberList(varName)
		if index == -1 {
			index = _func.GetMemberList().Len()
			_func.GetMemberList().PushBack(varName)
		}

		if m.WithOption() {
			_func.AddInstructionABC(opcode.StoreField, opcode.Memory, m.GetLeft().(ast.Expression).Compile(f, nil),
				compiler.NewSmallIntOperand(-int16(index+1)), r)
		} else {
			_func.AddInstructionABC(opcode.StoreField, opcode.Memory, m.GetLeft().(ast.Expression).Compile(f, nil),
				compiler.NewSmallIntOperand(int16(index)), r)
		}
	} else {
		index := _func.GetIndexOfLocalList(varName)
		if index != -1 && _func.CheckLocalVar(varName) {
			_func.AddInstructionABx(opcode.Move, opcode.Memory, compiler.NewRegisterOperand(_func.GetRegisterByLocalIndex(index)), r)
		} else if index = _func.GetIndexOfArgumentList(varName); index != -1 {
			_func.AddInstructionABx(opcode.Move, opcode.Memory, compiler.NewRegisterOperand(_func.GetRegisterByArgIndex(index)), r)
		} else if index = _func.GetIndexOfRefList(varName); index != -1 {
			_func.AddInstructionABx(opcode.Move, opcode.Memory, compiler.NewRefOperand(int16(index)), r)
		} else {
			index = _func.GetIndexOfRefList(varName)
			if index < 0 {
				index = _func.GetRefList().Len()
				_func.GetRefList().PushBack(varName)
			}
			_func.AddInstructionABx(opcode.Move, opcode.Memory, compiler.NewRefOperand(int16(index)), r)
		}
	}
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
	_func := f.(compiler.Function)

	switch impl.opType {
	case token.TokenTypeASSIGN:
		switch right := impl.right.(type) {
		case expression.Binary:
			if impl.left.(expression.Member).GetLeft() != nil {
				rr := right.(ast.Expression).Compile(f, nil)
				impl.compileStore(f, impl.left, rr)
				return rr
			} else {
				lr := impl.left.(ast.Expression).Compile(f, nil)
				return right.Compile(f, lr)
			}
		default:
			rr := right.(ast.Expression).Compile(f, r)
			impl.compileStore(f, impl.left, rr)
			return rr
		}
	case token.TokenTypeNULLISH:
		rv := compiler.NewRegisterOperand(_func.AllocRegister(""))
		impl.left.(ast.Expression).Compile(f, rv)
		jmp := _func.AddInstructionABx(opcode.JumpNull, opcode.Flow, rv, compiler.NewIntOperand(0))
		impl.right.(ast.Expression).Compile(f, rv)
		end := f.(compiler.Function).AddInstructionABx(opcode.Nop, opcode.Nop, compiler.NewSmallIntOperand(-1),
			compiler.NewIntOperand(0))
		jmp.Value.(*ast.Instruction).GetABx().B = end.Value.(*ast.Instruction).Index
		if r != nil {
			f.(compiler.Function).AddInstructionABx(opcode.Move, opcode.Memory, r, rv)
			return r
		}
		return rv
	default:
		lr := impl.left.(ast.Expression).Compile(f, nil)

		if r == nil {
			r = compiler.NewRegisterOperand(_func.AllocRegister(""))
		}

		switch impl.opType {
		case token.TokenTypeADD, token.TokenTypeADDASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Add, opcode.Math, r, lr, rr)
		case token.TokenTypeSUB, token.TokenTypeSUBASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Sub, opcode.Math, r, lr, rr)
		case token.TokenTypeMUL, token.TokenTypeMULASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Mul, opcode.Math, r, lr, rr)
		case token.TokenTypeDIV, token.TokenTypeDIVASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Div, opcode.Math, r, lr, rr)
		case token.TokenTypeREM, token.TokenTypeREMASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Rem, opcode.Math, r, lr, rr)
		case token.TokenTypeOR, token.TokenTypeORASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Or, opcode.Bit, r, lr, rr)
		case token.TokenTypeAND, token.TokenTypeANDASSIGN:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.And, opcode.Bit, r, lr, rr)
		case token.TokenTypeLESS:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Less, opcode.Logic, r, lr, rr)
		case token.TokenTypeLEQ:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.LessOrEqual, opcode.Logic, r, lr, rr)
		case token.TokenTypeGREATER:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Great, opcode.Logic, r, lr, rr)
		case token.TokenTypeGEQ:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.GreateOrEqual, opcode.Logic, r, lr, rr)
		case token.TokenTypeEQL:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.Equal, opcode.Logic, r, lr, rr)
		case token.TokenTypeNEQ:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.NotEqual, opcode.Logic, r, lr, rr)
		case token.TokenTypeLAND:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.LogicAnd, opcode.Logic, r, lr, rr)
		case token.TokenTypeLOR:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.LogicOr, opcode.Logic, r, lr, rr)
		case token.TokenTypeSHL:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.ShiftLeft, opcode.Math, r, lr, rr)
		case token.TokenTypeSHR:
			rr := impl.right.(ast.Expression).Compile(f, nil)
			_func.AddInstructionABC(opcode.ShiftRight, opcode.Math, r, lr, rr)
		}

		if impl.opType.WithAssign() {
			impl.compileStore(f, impl.left, r)
		}
	}

	return r
}

var _ expression.Binary = &Component{}

func NewBinary(owner, left, right interface{}, opType token.TokenType) *Component {
	return &Component{script.MakeComponentType(owner),
		left,
		right,
		opType,
		false,
	}
}
