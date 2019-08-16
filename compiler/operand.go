package compiler

import (
    "tklibs/script"
    "tklibs/script/opcode"
)

type Operand struct {
    *Register
    I int16
    script.Int
    script.Float
    Type    uint8
    IsSmall bool
}

func NewRegisterOperand(reg *Register) *Operand {
    return &Operand{
        Register: reg,
        Type:     opcode.Register,
    }
}

func NewRefOperand(refIndex int16) *Operand {
    return &Operand{
        I:    refIndex,
        Type: opcode.Reference,
    }
}

func NewSmallIntOperand(val int16) *Operand {
    return &Operand{
        I:       val,
        IsSmall: true,
        Type:    opcode.Integer,
    }
}

func NewIntOperand(val script.Int) *Operand {
    return &Operand{
        Int:  val,
        Type: opcode.Integer,
    }
}

func NewFloatOperand(val script.Float) *Operand {
    return &Operand{
        Float: val,
        Type:  opcode.None,
    }
}
