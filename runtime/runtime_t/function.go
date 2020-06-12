package runtime_t

import (
	"tklibs/script/debug"
	"tklibs/script/instruction"
)

type Function interface {
	GetInstructionList() []instruction.Instruction
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
	CopyFrom(src Function)
	GetAssembly() interface{}
	RegisterFunction(f uintptr)
	UnregisterFunction(f uintptr)
}

type NativeFunction interface {
	NativeCall(interface{}, interface{}, ...interface{}) interface{}
}
