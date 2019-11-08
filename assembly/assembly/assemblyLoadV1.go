package assembly

import (
    "bufio"
    "encoding/binary"

    "tklibs/script/runtime/function"
    "tklibs/script/runtime/runtime_t"
)

type assemblyFileHeaderV1 struct {
    *assemblyFileHeader
    functionCount uint32
}

type functionDeclV1 struct {
    instructionCount    uint32
    debugInfoCount      uint32
    localVars           []string
    arguments           []string
    refVars             []string
    members             []string
    nameIndex           uint32
    sourceFileNameCount uint32
    isScope             bool
    sourceFileNames     []string
    captureThis         bool
    maxRegisterCount    uint32
}

func loadV1(fh *assemblyFileHeader, impl *Component, reader *bufio.Reader) {
    fileHeader := &assemblyFileHeaderV1{
        assemblyFileHeader: fh,
    }

    binary.Read(reader, binary.LittleEndian, &fileHeader.functionCount)

    impl.intConstPool.Read(reader)
    impl.floatConstPool.Read(reader)
    impl.stringConstPool.Read(reader)

    functionDecls := make([]functionDeclV1, fileHeader.functionCount)
    for i := 0; i < int(fileHeader.functionCount); i++ {
        fd := &functionDecls[i]
        binary.Read(reader, binary.LittleEndian, &fd.instructionCount)
        binary.Read(reader, binary.LittleEndian, &fd.debugInfoCount)
        binary.Read(reader, binary.LittleEndian, &fd.nameIndex)
        binary.Read(reader, binary.LittleEndian, &fd.sourceFileNameCount)
        binary.Read(reader, binary.LittleEndian, &fd.isScope)
        binary.Read(reader, binary.LittleEndian, &fd.captureThis)
        binary.Read(reader, binary.LittleEndian, &fd.maxRegisterCount)

        fd.sourceFileNames = make([]string, fd.sourceFileNameCount)

        value := uint32(0)

        localCount := uint32(0)
        binary.Read(reader, binary.LittleEndian, &localCount)
        fd.localVars = make([]string, localCount)
        for j := uint32(0); j < localCount; j++ {
            binary.Read(reader, binary.LittleEndian, &value)
            fd.localVars[j] = impl.stringConstPool.Get(int(value)).(string)
        }

        argCount := uint8(0)
        binary.Read(reader, binary.LittleEndian, &argCount)
        fd.arguments = make([]string, argCount)
        for j := uint8(0); j < argCount; j++ {
            binary.Read(reader, binary.LittleEndian, &value)
            fd.arguments[j] = impl.stringConstPool.Get(int(value)).(string)
        }

        refCount := uint32(0)
        binary.Read(reader, binary.LittleEndian, &refCount)
        fd.refVars = make([]string, refCount)
        for j := uint32(0); j < refCount; j++ {
            binary.Read(reader, binary.LittleEndian, &value)
            fd.refVars[j] = impl.stringConstPool.Get(int(value)).(string)
        }

        memberCount := uint32(0)
        binary.Read(reader, binary.LittleEndian, &memberCount)
        fd.members = make([]string, memberCount)
        for j := uint32(0); j < memberCount; j++ {
            binary.Read(reader, binary.LittleEndian, &value)
            fd.members[j] = impl.stringConstPool.Get(int(value)).(string)
        }

        for i := range fd.sourceFileNames {
            binary.Read(reader, binary.LittleEndian, &value)
            fd.sourceFileNames[i] = impl.stringConstPool.Get(int(value)).(string)
        }
    }

    impl.functions = make([]interface{}, fileHeader.functionCount)

    for i := 0; i < int(fileHeader.functionCount); i++ {
        f := &struct {
            *function.Component
        }{}

        fd := functionDecls[i]
        f.Component = function.NewFunctionComponent(f,
            int(fd.instructionCount),
            int(fd.debugInfoCount),
            impl.stringConstPool.Get(int(fd.nameIndex)).(string),
            fd.sourceFileNames,
            fd.localVars,
            fd.arguments,
            fd.refVars,
            fd.members,
            fd.isScope,
            fd.captureThis,
            fd.maxRegisterCount,
        )

        impl.functions[i] = f
    }

    for i := 0; i < int(fileHeader.functionCount); i++ {
        instList := impl.functions[i].(runtime_t.Function).GetInstructionList()
        debugInfoList := impl.functions[i].(runtime_t.Function).GetDebugInfoList()

        for j := 0; j < len(instList); j++ {
            binary.Read(reader, binary.LittleEndian, &instList[j])
        }

        for j := 0; j < len(debugInfoList); j++ {
            binary.Read(reader, binary.LittleEndian, &debugInfoList[j].Line)
            binary.Read(reader, binary.LittleEndian, &debugInfoList[j].PC)
            binary.Read(reader, binary.LittleEndian, &debugInfoList[j].SourceIndex)
        }
    }
}
