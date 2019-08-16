package assembly

import (
    "bufio"
    "encoding/binary"

    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/debug"
)

func saveV1(fh *assemblyFileHeader, impl *Component, writer *bufio.Writer) {
    fileHeader := &assemblyFileHeaderV1{
        assemblyFileHeader: fh,
    }
    fileHeader.functionCount = uint32(len(impl.functions))
    binary.Write(writer, binary.LittleEndian, fileHeader.functionCount)

    impl.intConstPool.Write(writer)
    impl.floatConstPool.Write(writer)
    impl.stringConstPool.Write(writer)

    functionDecls := make([]functionDeclV1, fileHeader.functionCount)
    for i := 0; i < len(impl.functions); i++ {
        _func := impl.functions[i].(compiler.Function)

        fd := functionDecls[i]

        fd.instructionCount = uint32(_func.GetInstructionList().Len())
        fd.debugInfoCount = uint32(_func.GetDebugInfoList().Len())
        fd.nameIndex = *_func.GetNameIndexPointer()
        fd.sourceFileNameCount = uint32(len(_func.GetSourceNames()))
        fd.sourceFileNames = _func.GetSourceNames()
        fd.isScope = _func.IsScope()
        fd.captureThis = _func.GetCaptureThis()
        fd.maxRegisterCount = uint32(_func.GetMaxRegisterCount())

        binary.Write(writer, binary.LittleEndian, fd.instructionCount)
        binary.Write(writer, binary.LittleEndian, fd.debugInfoCount)
        binary.Write(writer, binary.LittleEndian, fd.nameIndex)
        binary.Write(writer, binary.LittleEndian, fd.sourceFileNameCount)
        binary.Write(writer, binary.LittleEndian, fd.isScope)
        binary.Write(writer, binary.LittleEndian, fd.captureThis)
        binary.Write(writer, binary.LittleEndian, fd.maxRegisterCount)

        binary.Write(writer, binary.LittleEndian, uint32(_func.GetLocalList().Len()))
        for it := _func.GetLocalList().Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, uint32(it.Value.(int)))
        }

        binary.Write(writer, binary.LittleEndian, uint8(_func.GetArgList().Len()))
        for it := _func.GetArgList().Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, uint32(it.Value.(int)))
        }

        binary.Write(writer, binary.LittleEndian, uint32(_func.GetRefList().Len()))
        for it := _func.GetRefList().Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, uint32(it.Value.(int)))
        }

        binary.Write(writer, binary.LittleEndian, uint32(_func.GetMemberList().Len()))
        for it := _func.GetMemberList().Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, uint32(it.Value.(int)))
        }

        for _, sourceName := range fd.sourceFileNames {
            binary.Write(writer, binary.LittleEndian, uint32(impl.stringConstPool.Insert(sourceName)))
        }
    }

    for i := 0; i < len(impl.functions); i++ {
        _func := impl.functions[i].(compiler.Function)
        instList := _func.GetInstructionList()
        debugInfoList := _func.GetDebugInfoList()

        for it := instList.Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, &it.Value.(*ast.Instruction).Instruction)
        }

        for it := debugInfoList.Front(); it != nil; it = it.Next() {
            binary.Write(writer, binary.LittleEndian, &it.Value.(*debug.Info).Line)
            binary.Write(writer, binary.LittleEndian, &it.Value.(*debug.Info).PC)
            binary.Write(writer, binary.LittleEndian, &it.Value.(*debug.Info).SourceIndex)
        }
    }
}
