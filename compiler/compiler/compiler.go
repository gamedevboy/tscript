package compiler

import (
    "container/list"
    "fmt"
    "sort"

    "tklibs/script"
    "tklibs/script/assembly"
    assemblyImpl "tklibs/script/assembly/assembly"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/lexer"
    "tklibs/script/compiler/parser"
    parserComponent "tklibs/script/compiler/parser/parser"
)

type Component struct {
    script.ComponentType
    fileList    list.List
    sourceList list.List
    funcMap    map[uint32]interface{}
    funcStack  list.List
    funcIndex  uint32
}

func (impl *Component) AddFile(fileName string) {
    impl.fileList.PushBack(fileName)
}

func (impl *Component) AddSource(code string) {
    impl.sourceList.PushBack(code)
}

func (impl *Component) Compile() (interface{}, *list.List, error) {
    l := &struct {
        *lexer.Component
    }{}
    l.Component = lexer.NewLexer(l)

    p := parserComponent.NewParser()

    tokenList := list.New()

    for it := impl.fileList.Front(); it != nil; it = it.Next() {
        fileName := it.Value.(string)

        tl, err := l.ParseFile(fileName, tokenList)
        if err != nil {
            return nil, tokenList, err
        }

        tokenList = tl
    }

    for it := impl.sourceList.Front(); it != nil; it = it.Next() {
        tokenList = l.ParseFromRunes("[SOURCE]", []rune(it.Value.(string)), tokenList)
    }

    asm := &struct {
        *assemblyImpl.Component
    }{}

    entryFunction := newCompilerFunction(asm)

    bs := &struct {
        *block.Component
    }{}
    bs.Component = block.NewBlock(bs)
    ef := entryFunction.(compiler.Function)
    ef.SetName("global")
    ef.SetBlockStatement(bs)
    ef.SetScope(true)
    tokenStart := tokenList.Front()

    p.(parser.BlockParser).ParseBlock(bs, tokenStart)

    impl.funcIndex = 0
    impl.funcStack.PushBack(entryFunction)
    impl.funcMap[0] = entryFunction
    impl.visitForFunctionScan(bs, asm)

    functions := make([]interface{}, len(impl.funcMap))

    asm.Component = assemblyImpl.NewScriptAssemblyWithFunctions(asm, functions)

    for _, f := range impl.funcMap {
        functions[*f.(compiler.Function).GetFunctionIndexPointer()] = f
    }

    for i, f := range functions {
        functions[i] = impl.compile(f)
    }

    return asm, tokenList, nil
}

func (impl *Component) addFunc(f interface{}) {
    _func := f.(compiler.Function)
    impl.funcMap[*_func.GetFunctionIndexPointer()] = f
}

func (impl *Component) visitForFunctionScan(astNode, asm interface{}) {
    if astNode == nil {
        return
    }

    curFunc := impl.funcStack.Back().Value

    switch target := astNode.(type) {
    case statement.Block:
        statementList := target.GetStatementList()
        for it := statementList.Front(); it != nil; it = it.Next() {
            impl.visitForFunctionScan(it.Value, asm)
        }
    case expression.Call:
        impl.visitForFunctionScan(target.GetExpression(), asm)
        impl.visitForFunctionScan(target.GetArgList(), asm)
    case statement.If:
        impl.visitForFunctionScan(target.GetCondition(), asm)
        impl.visitForFunctionScan(target.GetBody(), asm)
        impl.visitForFunctionScan(target.GetElseBody(), asm)
    case statement.For:
        impl.visitForFunctionScan(target.GetInit(), asm)
        impl.visitForFunctionScan(target.GetCondition(), asm)
        impl.visitForFunctionScan(target.GetStep(), asm)
        impl.visitForFunctionScan(target.GetBody(), asm)
    case statement.While:
        impl.visitForFunctionScan(target.GetCondition(), asm)
        impl.visitForFunctionScan(target.GetBody(), asm)
    case statement.Switch:
        for it := target.GetCaseList().Front(); it != nil; it = it.Next() {
            impl.visitForFunctionScan(it.Value.(statement.Case).GetBlock(), asm)
        }
        impl.visitForFunctionScan(target.GetDefaultCase(), asm)
    case statement.Decl:
        if impl.funcStack.Len() == 1 {
            target.SetGlobal(true)
        }

        if !target.IsGlobal() {
            curFunc.(compiler.Function).GetLocalList().PushBack(target.GetName())
        }

        impl.visitForFunctionScan(target.GetExpression(), asm)
    case statement.Return:
        impl.visitForFunctionScan(target.GetExpression(), asm)
    case expression.Object:
        keys := make([]string, len(target.GetKeyValueMap()))

        idx := 0
        for name := range target.GetKeyValueMap() {
            keys[idx] = name
            idx++
        }

        sort.Strings(keys)

        for _, name := range keys {
            impl.visitForFunctionScan(target.GetKeyValueMap()[name], asm)
        }
    case expression.Array:
        impl.visitForFunctionScan(target.GetArgListExpression(), asm)
    case expression.ArgList:
        for it := target.GetExpressionList().Front(); it != nil; it = it.Next() {
            impl.visitForFunctionScan(it.Value, asm)
        }
    case expression.Function:
        impl.funcIndex++
        _f := newCompilerFunction(asm)
        _func := _f.(compiler.Function)
        _func.SetCaptureThis(target.GetCaptureThis())
        _func.SetBlockStatement(target.GetBlock())
        //debugInfo := target.(debug.Info)
        *_func.GetFunctionIndexPointer() = impl.funcIndex
        switch target.GetName() {
        case "":
            if impl.funcStack.Len() > 0 {
                topFuncName := impl.funcStack.Back().Value.(compiler.Function).GetName()
                //line := debugInfo.GetLine()
                _func.SetName(fmt.Sprintf("%v.%v", topFuncName, impl.funcIndex))
            }
        default:
            _func.SetName(target.GetName())
        }

        target.SetMetaIndex(impl.funcIndex)
        argList := target.GetArgList()
        if argList != nil {
            for it := argList.(expression.ArgList).GetExpressionList().Front();
                it != nil; it = it.Next() {
                m := it.Value.(expression.Member)
                _func.GetArgList().PushBack(m.GetRight().(string))
            }
        }
        impl.addFunc(_f)
        if impl.funcStack.Len() > 0 {
            impl.funcStack.Back().Value.(compiler.Function).SetScope(true)
        }
        impl.funcStack.PushBack(_f)
        impl.visitForFunctionScan(target.GetBlock(), asm)
        impl.funcStack.Remove(impl.funcStack.Back())
    case expression.Binary:
        impl.visitForFunctionScan(target.GetLeft(), asm)
        impl.visitForFunctionScan(target.GetRight(), asm)
    case expression.Const:
    case statement.Expression:
        impl.visitForFunctionScan(target.GetExpression(), asm)
    }
}

func insertIntoStringPool(stringPool assembly.ConstPool, stringList *list.List) {
    for it := stringList.Front(); it != nil; it = it.Next() {
        it.Value = stringPool.Insert(it.Value)
    }
}

func (impl *Component) compile(v interface{}) interface{} {
    _func := v.(compiler.Function)
    _func.GetBlockStatement().(ast.Statement).Compile(v)

    constStringPool := _func.GetAssembly().(script.Assembly).GetStringConstPool().(assembly.ConstPool)
    *_func.GetNameIndexPointer() = uint32(constStringPool.Insert(_func.GetName()))

    insertIntoStringPool(constStringPool, _func.GetLocalList())
    insertIntoStringPool(constStringPool, _func.GetArgList())
    insertIntoStringPool(constStringPool, _func.GetRefList())
    insertIntoStringPool(constStringPool, _func.GetMemberList())

    return v
}

func NewCompiler(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        funcMap:       make(map[uint32]interface{}),
    }
}
