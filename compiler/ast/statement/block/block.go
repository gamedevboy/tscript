package block

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/debug"
    debug2 "tklibs/script/debug"
)

type Component struct {
    debug.Component
    script.ComponentType
    statementList list.List
}

func (impl *Component) GetStatementList() *list.List {
    return &impl.statementList
}

func (impl *Component) Compile(f interface{}) *list.Element {
    _func := f.(compiler.Function)

    _func.PushBlock()
    defer _func.PopBlock()

    var start *list.Element
    stringConstPool := _func.GetAssembly().(script.Assembly).GetStringConstPool()
    
    for it := impl.statementList.Front(); it != nil; it = it.Next() {
        statement := it.Value.(ast.Statement)
        debugInfo := statement.(debug.Info)

        stringConstPool.Insert(debugInfo.GetFilePath())

        _func.GetDebugInfoList().PushBack(&debug2.Info{
            Line:        uint32(debugInfo.GetLine()),
            PC:          uint32(_func.GetInstructionList().Len()),
            SourceIndex: uint32(_func.AddSourceFile(debugInfo.GetFilePath())),
        })

        if p := statement.Compile(f); start == nil {
            start = p
        }

        _func.ReleaseAllRegisters()
    }

    return start
}

func NewBlock(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}
