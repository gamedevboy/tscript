package runtime

import (
    "tklibs/script/debug"
    "tklibs/script/instruction"
)

type Function interface {
    GetInstructionList() []instruction.Instruction
    SetInstructionList([]instruction.Instruction)
    GetDebugInfoList() []debug.Info
    DumpString() string
    GetArguments() []string
    GetLocalVars() []string
    GetRefVars() []string
    GetMembers() []string
    GetName() string
    IsScope() bool
    GetSourceNames() []string
    IsCaptureThis() bool
    GetMaxRegisterCount() int
}
