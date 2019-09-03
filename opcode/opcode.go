package opcode

import "tklibs/script"

type OperandABC struct {
    Code, Type uint8
    A, B, C    int16
}

type OperandABx struct {
    Code, Type uint8
    A          int16
    B          script.Int
}

type OperandABm struct {
    Code, Type uint8
    A          int16
    B          script.Float
}

const (
    None uint8 = iota
    Register
    Reference
    Integer
)

const (
    Nop uint8 = iota
    Memory
    Math
    Bit
    Logic
    Flow
    Const
)

// Memory
const (
    Move uint8 = iota
    LoadField
    StoreField
    LoadElement
    StoreElement
    Object
    Array
)

// Const
const (
    Load uint8 = iota
    LoadFunc
    LoadNil
)

// Math
const (
    Add uint8 = iota
    Sub
    Mul
    Div
    Neg
    Inc
    Dec
    Rem
    Shift
)

// Logic
const (
    Less uint8 = iota
    Great
    Equal
    NotEqual
    LessOrEqual
    GreateOrEqual
    LogicAnd
    LogicOr
    LogicNot
)

// Bit
const (
    And uint8 = iota
    Or
    Not
    Xor
)

// Flow
const (
    Jump uint8 = iota
    JumpWhenFalse
    Call
    NewCall
    Ret
)

// Const Type
const (
    ConstInt64 script.Int = iota
    ConstFloat64
    ConstString
)
