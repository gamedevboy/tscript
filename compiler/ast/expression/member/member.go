package member

import (
	"fmt"
	"strings"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/opcode"
)

type Component struct {
	script.ComponentType
	left, right interface{}
	option      bool
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	if impl.left != nil {
		impl.left.(ast.Node).Format(ident, formatBuilder)
		switch v := impl.right.(type) {
		case string:
			formatBuilder.WriteString(".")
			formatBuilder.WriteString(v)
		default:
			formatBuilder.WriteString("[")
			v.(ast.Node).Format(ident, formatBuilder)
			formatBuilder.WriteString("]")
		}
	} else {
		formatBuilder.WriteString(impl.right.(string))
	}
}

func (impl *Component) WithOption() bool {
	return impl.option
}

var _ expression.Member = &Component{}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
	_func := f.(compiler.Function)

	if impl.left != nil {
		if r == nil {
			r = compiler.NewRegisterOperand(_func.AllocRegister(""))
		}
		lr := impl.left.(ast.Expression).Compile(f, nil)

		switch varName := impl.right.(type) {
		case string:
			impl.compileLoad(_func, varName, r, lr)
		case expression.Const:
			switch vn := varName.GetValue().(type) {
			case script.String:
				impl.compileLoad(_func, string(vn), r, lr)
			default:
				_func.AddInstructionABC(opcode.LoadElement, opcode.Memory, r, lr, varName.Compile(f, nil))
			}
		case ast.Expression:
			_func.AddInstructionABC(opcode.LoadElement, opcode.Memory, r, lr, varName.Compile(f, nil))
		}
		return r
	} else {
		if varName, ok := impl.right.(string); ok {
			switch varName {
			case "this":
				if r == nil {
					return compiler.NewRegisterOperand(&compiler.Register{Index: 1})
				} else {
					_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRegisterOperand(&compiler.Register{Index: 1}))
					return r
				}
			case "null":
				if r == nil {
					n := compiler.NewRegisterOperand(_func.AllocRegister(""))
					_func.AddInstructionABx(opcode.LoadNil, opcode.Const, n, compiler.NewIntOperand(0))
					return n
				} else {
					_func.AddInstructionABx(opcode.LoadNil, opcode.Const, r, compiler.NewIntOperand(0))
					return r
				}
			default:
				index := _func.GetIndexOfLocalList(varName)
				if index != -1 && _func.CheckLocalVar(varName) {
					if r == nil {
						return compiler.NewRegisterOperand(_func.GetRegisterByLocalIndex(index))
					} else {
						_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRegisterOperand(_func.GetRegisterByLocalIndex(index)))
						return r
					}
				} else if index = _func.GetIndexOfArgumentList(varName); index != -1 {
					if r == nil {
						return compiler.NewRegisterOperand(_func.GetRegisterByArgIndex(index))
					} else {
						_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRegisterOperand(_func.GetRegisterByArgIndex(index)))
						return r
					}
				} else if index = _func.GetIndexOfRefList(varName); index != -1 {
					if r == nil {
						return compiler.NewRefOperand(int16(index))
					} else {
						_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRefOperand(int16(index)))
						return r
					}
				} else {
					index = _func.GetIndexOfRefList(varName)
					if index < 0 {
						index = _func.GetRefList().Len()
						_func.GetRefList().PushBack(varName)
					}

					if r == nil {
						return compiler.NewRefOperand(int16(index))
					} else {
						_func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRefOperand(int16(index)))
						return r
					}
				}
			}
		} else {
			// todo
			panic("")
		}
	}
}

func (impl *Component) compileLoad(_func compiler.Function, varName string, r *compiler.Operand, lr *compiler.Operand) {
	index := _func.GetIndexOfMemberList(varName)
	if index == -1 {
		memberList := _func.GetMemberList()
		index = memberList.Len()
		memberList.PushBack(varName)
	}
	if impl.option {
		_func.AddInstructionABC(opcode.LoadField, opcode.Memory, r, lr, compiler.NewSmallIntOperand(-(int16(index) + 1)))
	} else {
		_func.AddInstructionABC(opcode.LoadField, opcode.Memory, r, lr, compiler.NewSmallIntOperand(int16(index)))
	}
}

func (impl *Component) String() string {
	if impl.left != nil {
		return fmt.Sprint(impl.left, ".", impl.right)
	}

	return fmt.Sprint(impl.right)
}

func (impl *Component) GetLeft() interface{} {
	return impl.left
}

func (impl *Component) GetRight() interface{} {
	return impl.right
}

func NewMember(owner, left, right interface{}, option bool) *Component {
	return &Component{script.MakeComponentType(owner),
		left,
		right,
		option,
	}
}
