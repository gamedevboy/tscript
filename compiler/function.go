package compiler

import (
    "container/list"
)

type Function interface {
    GetFunctionIndexPointer() *uint32
    GetLocalList() *list.List
    GetArgList() *list.List
    GetRefList() *list.List
    GetMemberList() *list.List
    GetNameIndexPointer() *uint32
    GetSourceNames() []string
    GetName() string
    SetName(name string)
    SetBlockStatement(interface{})
    GetBlockStatement() interface{}
    GetIndexOfLocalList(string) int
    GetIndexOfArgumentList(string) int
    GetIndexOfRefList(string) int
    GetIndexOfMemberList(string) int
    GetInstructionList() *list.List
    GetDebugInfoList() *list.List
    AddInstructionABC(code, _type uint8, a, b, c *Operand) *list.Element
    AddInstructionABx(code, _type uint8, a, b *Operand) *list.Element
    AddInstructionABm(code, _type uint8, a, b *Operand) *list.Element
    GetAssembly() interface{}
    PushBreakList()
    PopBreakList()
    GetBreakList() *list.List
    GetContinueList() *list.List
    PushBlock()
    PopBlock()
    CheckLocalVar(name string) bool
    SetScope(bool)
    IsScope() bool
    GetRegisterByLocalIndex(index int) *Register
    GetRegisterByArgIndex(index int) *Register
    ReleaseAllRegisters()
    AllocRegister(tag string) *Register
    GetRegisterCount() int16
    AddSourceFile(filePath string) int
    SetCaptureThis(val bool)
    GetCaptureThis() bool
    PushRegisters()
    PopRegisters()
    GetMaxRegisterCount() int

}
