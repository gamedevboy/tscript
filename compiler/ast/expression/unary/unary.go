package unary

import (
    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/expression/binary"
    _const "tklibs/script/compiler/ast/expression/const"
    "tklibs/script/compiler/token"
    "tklibs/script/opcode"
)

type Component struct {
    script.ComponentType
    expression interface{}
    tokenType  token.TokenType
}

func (c *Component) String() string {
    panic("implement me")
}

func (c *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
    _func := f.(compiler.Function)
    switch c.tokenType {
    case token.TokenTypeINC:
        switch target := c.expression.(type) {
        case expression.Member:
            c := &struct {
                *_const.Component
            }{}
            c.Component = _const.NewConst(c, script.Int(1))
            b := &struct {
                *binary.Component
            }{}
            b.Component = binary.NewBinary(b, target, c, token.TokenTypeADDASSIGN)
            return b.Compile(f, r)
        case ast.Expression:
            _func.AddInstructionABx(opcode.Inc, opcode.Math, target.Compile(f, r), nil)
        }
    case token.TokenTypeDEC:
        switch target := c.expression.(type) {
        case expression.Member:
            c := &struct {
                *_const.Component
            }{}
            c.Component = _const.NewConst(c, script.Int(1))
            b := &struct {
                *binary.Component
            }{}
            b.Component = binary.NewBinary(b, target, c, token.TokenTypeSUBASSIGN)
            return b.Compile(f, r)
        case ast.Expression:
            _func.AddInstructionABx(opcode.Inc, opcode.Math, target.Compile(f, r), nil)
        }
    case token.TokenTypeSUB:
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(""))
        }
        _func.AddInstructionABx(opcode.Neg, opcode.Math, r, c.expression.(ast.Expression).Compile(f, nil))
    case token.TokenTypeLNOT:
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(""))
        }
        _func.AddInstructionABx(opcode.LogicNot, opcode.Logic, r, c.expression.(ast.Expression).Compile(f, nil))
    default:
        panic("Not support")
    }

    return r
}

var _ expression.Unary = &Component{}

func NewUnary(owner, expression interface{}, tokenType token.TokenType) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        expression:    expression,
        tokenType:     tokenType,
    }
}
