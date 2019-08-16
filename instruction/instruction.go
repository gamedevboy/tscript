package instruction

import (
    "unsafe"

    "tklibs/script/opcode"
)

type Instruction struct {
    opcode.OperandABC
}

func (impl *Instruction) GetABC() *opcode.OperandABC {
    return &impl.OperandABC
}

func (impl *Instruction) GetABx() *opcode.OperandABx {
    return (*opcode.OperandABx)(unsafe.Pointer(impl))
}

func (impl *Instruction) GetABm() *opcode.OperandABm {
    return (*opcode.OperandABm)(unsafe.Pointer(impl))
}
