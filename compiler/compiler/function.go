package compiler

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/opcode"
)

type compilerFunction struct {
    script.ComponentType
    blockStatement   interface{}
    argList          list.List
    localVarList     list.List
    refList          list.List
    memberList       list.List
    metaIndex        uint32
    nameIndex        uint32
    sourceFileNames  []string
    instructionList  list.List
    debugInfoList    list.List
    asm              interface{}
    name             string
    breakList        list.List
    continueList     list.List
    isScope          bool
    allocatedRegList *list.List
    sourceFile       string
    captureThis      bool
    registerList     list.List
    maxRegisterCount int
}

func (impl *compilerFunction) GetContinueList() *list.List {
    return &impl.continueList
}

func (impl *compilerFunction) GetMaxRegisterCount() int {
    return impl.maxRegisterCount
}

func (impl *compilerFunction) SetCaptureThis(val bool) {
    impl.captureThis = val
}

func (impl *compilerFunction) GetCaptureThis() bool {
    return impl.captureThis
}

var _ compiler.Function = &compilerFunction{}

func (impl *compilerFunction) GetSourceNames() []string {
    return impl.sourceFileNames
}

func (impl *compilerFunction) AddSourceFile(filePath string) int {
    for i, f := range impl.sourceFileNames {
        if f == filePath {
            return i
        }
    }

    impl.sourceFileNames = append(impl.sourceFileNames, filePath)
    return len(impl.sourceFileNames) - 1
}

func (impl *compilerFunction) GetRegisterCount() int16 {
    return int16(2 + impl.allocatedRegList.Len() + impl.argList.Len() + impl.localVarList.Len())
}

func (impl *compilerFunction) GetRegisterByArgIndex(index int) *compiler.Register {
    return &compiler.Register{
        Index: int16(2 + index),
        Tag:   fmt.Sprintf("Arg:%v", 2+index),
    }
}

func (impl *compilerFunction) GetRegisterByLocalIndex(index int) *compiler.Register {
    return &compiler.Register{
        Index: int16(2 + impl.argList.Len() + index),
        Tag:   fmt.Sprintf("Local:%v", 2+impl.argList.Len()+index),
    }
}

func (impl *compilerFunction) ReleaseAllRegisters() {
    impl.allocatedRegList.Init()

    if impl.registerList.Len() > 0 {
        l := impl.registerList.Back().Value.(*list.List)
        for it := l.Front(); it != nil; it = it.Next() {
            impl.allocatedRegList.PushBack(it.Value)
        }
    }
}

func (impl *compilerFunction) PushRegisters() {
    l := list.New()

    for it := impl.allocatedRegList.Front(); it != nil; it = it.Next() {
        l.PushBack(it.Value)
    }

    impl.registerList.PushBack(l)
}

func (impl *compilerFunction) PopRegisters() {
    l := impl.registerList.Back().Value.(*list.List)
    impl.allocatedRegList.Init()
    for it := l.Front(); it != nil; it = it.Next() {
        impl.allocatedRegList.PushBack(it.Value)
    }
    impl.registerList.Remove(impl.registerList.Back())
}

func (impl *compilerFunction) AllocRegister(tag string) *compiler.Register {
    if impl.allocatedRegList.Len() > 0 {
        for it := impl.allocatedRegList.Front(); it != nil; it = it.Next() {
            reg := it.Value.(*compiler.Register)
            if reg.Tag == tag && tag != "" {
                return reg
            }
        }
    }

    ret := &compiler.Register{
        Index: int16(impl.allocatedRegList.Len() + 2 + impl.GetLocalList().Len() + impl.GetArgList().Len()),
        Tag:   tag,
    }

    impl.allocatedRegList.PushBack(ret)

    if impl.allocatedRegList.Len() > impl.maxRegisterCount {
        impl.maxRegisterCount = impl.allocatedRegList.Len()
    }

    return ret
}

func (impl *compilerFunction) AddInstructionABC(code, _type uint8, a, b, c *compiler.Operand) *list.Element {
    inst := &ast.Instruction{}
    inst.Index = impl.instructionList.Len()
    inst.Code = code
    inst.Type = _type << 4

    if b != nil {
        inst.Type |= b.Type << 2
        switch b.Type {
        case opcode.Register:
            inst.B = b.Index
        default:
            inst.B = b.I
        }
    }
    if c != nil {
        inst.Type |= c.Type & 3
        switch c.Type {
        case opcode.Register:
            inst.C = c.Index
        default:
            inst.C = c.I
        }
    }

    if a.Type == opcode.Register {
        inst.A = a.Index
    } else {
        inst.A = -a.I - 1
    }

    return impl.instructionList.PushBack(inst)
}

func (impl *compilerFunction) AddInstructionABx(code, _type uint8, a, b *compiler.Operand) *list.Element {
    inst := &ast.Instruction{}
    inst.Index = impl.instructionList.Len()
    inst.Code = code
    inst.Type = _type << 4
    ptr := inst.GetABx()

    if b != nil {
        inst.Type |= b.Type << 2
        switch b.Type {
        case opcode.Register:
            ptr.B = script.Int(b.Index)
        case opcode.Reference:
            ptr.B = script.Int(b.I)
        default:
            if b.IsSmall {
                ptr.B = script.Int(b.I)
            } else {
                ptr.B = b.Int
            }
        }
    }
    if a.Type == opcode.Register {
        ptr.A = a.Index
    } else {
        ptr.A = -a.I - 1
    }

    return impl.instructionList.PushBack(inst)
}

func (impl *compilerFunction) AddInstructionABm(code, _type uint8, a, b *compiler.Operand) *list.Element {
    inst := &ast.Instruction{}
    inst.Index = impl.instructionList.Len()
    inst.Code = code
    inst.Type = _type << 4
    ptr := inst.GetABm()

    if b != nil {
        inst.Type |= b.Type << 2
        ptr.B = b.Float
    }

    if a.Type == opcode.Register {
        ptr.A = a.Register.Index
    } else {
        ptr.A = -a.I - 1
    }

    return impl.instructionList.PushBack(inst)
}

func (impl *compilerFunction) SetScope(value bool) {
    impl.isScope = value
}

func (impl *compilerFunction) IsScope() bool {
    return impl.isScope
}

func (impl *compilerFunction) GetFunctionIndexPointer() *uint32 {
    return &impl.metaIndex
}

func (impl *compilerFunction) GetLocalList() *list.List {
    return &impl.localVarList
}

func (impl *compilerFunction) GetArgList() *list.List {
    return &impl.argList
}

func (impl *compilerFunction) GetRefList() *list.List {
    return &impl.refList
}

func (impl *compilerFunction) GetMemberList() *list.List {
    return &impl.memberList
}

func (impl *compilerFunction) SetBlockStatement(bs interface{}) {
    impl.blockStatement = bs
}

func (impl *compilerFunction) GetBlockStatement() interface{} {
    return impl.blockStatement
}

func (impl *compilerFunction) GetDebugInfoList() *list.List {
    return &impl.debugInfoList
}

func (impl *compilerFunction) GetIndexOfLocalList(value string) int {
    ret := 0

    for it := impl.localVarList.Front(); it != nil; it = it.Next() {
        if it.Value.(string) == value {
            return ret
        }
        ret++
    }

    return -1
}

func (impl *compilerFunction) GetIndexOfMemberList(value string) int {
    ret := 0

    for it := impl.memberList.Front(); it != nil; it = it.Next() {
        if it.Value.(string) == value {
            return ret
        }
        ret++
    }

    return -1
}

func (impl *compilerFunction) GetIndexOfArgumentList(value string) int {
    ret := 0

    for it := impl.argList.Front(); it != nil; it = it.Next() {
        if it.Value.(string) == value {
            return ret
        }
        ret++
    }

    return -1
}

func (impl *compilerFunction) GetIndexOfRefList(value string) int {
    ret := 0

    for it := impl.refList.Front(); it != nil; it = it.Next() {
        if it.Value.(string) == value {
            return ret
        }
        ret++
    }

    return -1
}

func (impl *compilerFunction) GetInstructionList() *list.List {
    return &impl.instructionList
}

func (impl *compilerFunction) GetAssembly() interface{} {
    return impl.asm
}

func (impl *compilerFunction) GetName() string {
    return impl.name
}

func (impl *compilerFunction) SetName(name string) {
    impl.name = name
}

func (impl *compilerFunction) GetNameIndexPointer() *uint32 {
    return &impl.nameIndex
}

func (impl *compilerFunction) GetBreakList() *list.List {
    return impl.breakList.Back().Value.(*list.List)
}

func (impl *compilerFunction) PushBreakList() {
    impl.breakList.PushBack(list.New())
}

func (impl *compilerFunction) PopBreakList() {
    impl.breakList.Remove(impl.breakList.Back())
}

func (impl *compilerFunction) PushBlock() {

}

func (impl *compilerFunction) PopBlock() {

}

func (impl *compilerFunction) CheckLocalVar(name string) bool {
    return true
}

func newCompilerFunction(asm interface{}) interface{} {
    _func := &struct {
        *compilerFunction
    }{}

    _func.compilerFunction = &compilerFunction{
        ComponentType:    script.MakeComponentType(_func),
        asm:              asm,
        nameIndex:        0xffffffff,
        allocatedRegList: list.New(),
        sourceFileNames:  make([]string, 0, 1),
    }

    return _func
}
