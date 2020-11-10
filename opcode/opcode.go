package opcode

type OperandABC struct {
    Code, Type uint8
    A, B, C    int16
}

type OperandABx struct {
    Code, Type uint8
    A          int16
    B          int32
}

type OperandABm struct {
    Code, Type uint8
    A          int16
    B          float32
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
    ShiftLeft
    ShiftRight
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
    JumpNull
    JumpTo
    Call
    Ret
)

// Const Type
const (
    ConstInt64 int = iota
    ConstFloat64
    ConstString
    ConstBool
)
