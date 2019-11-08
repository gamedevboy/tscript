package ast

import "tklibs/script/instruction"

type Instruction struct {
    instruction.Instruction
    Index int32
}
