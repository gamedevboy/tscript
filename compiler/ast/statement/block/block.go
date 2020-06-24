package block

import (
	"container/list"
	"strings"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/statement"
	"tklibs/script/compiler/debug"
	debug2 "tklibs/script/debug"
)

type Component struct {
	debug.Component
	script.ComponentType
	statementList list.List
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	first := impl.statementList.Front()
	if first != nil {
		if t, ok := first.Value.(statement.Trivia); ok && t.(debug.Info).GetLine() == impl.GetLine() && t.GetContent()[1] == '/' {
			formatBuilder.WriteString(" ")
		} else {
			formatBuilder.WriteString("\n")
		}
	}

	for it := impl.statementList.Front(); it != nil; it = it.Next() {
		st := it.Value.(ast.Statement)
		debugInfo := st.(debug.Info)
		skipLine := debugInfo.GetSkipLine()

		for skipLine > 0 {
			formatBuilder.WriteString("\n")
			skipLine--
		}

		prev := it.Prev()

		if prev != nil && prev.Value.(debug.Info).GetLine() != debugInfo.GetLine() {
			formatBuilder.WriteString(strings.Repeat(" ", ident))
		} else {
			if _, ok := it.Value.(statement.Trivia); ok {
				formatBuilder.WriteString(" ")
			}
		}

		st.Format(ident, formatBuilder)
		next := it.Next()

		if next == nil {
			formatBuilder.WriteString("\n")
			return
		}

		if t, ok := next.Value.(statement.Trivia); !(ok && t.(debug.Info).GetLine() == debugInfo.GetLine()) {
			formatBuilder.WriteString("\n")
		}
	}
}

var _ ast.Statement = &Component{}

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
		statementItem := it.Value.(ast.Statement)

		debugInfo := statementItem.(debug.Info)

		stringConstPool.Insert(debugInfo.GetFilePath())

		_func.GetDebugInfoList().PushBack(&debug2.Info{
			Line:        uint32(debugInfo.GetLine()),
			PC:          uint32(_func.GetInstructionList().Len()),
			SourceIndex: uint32(_func.AddSourceFile(debugInfo.GetFilePath())),
		})

		if p := statementItem.Compile(f); start == nil {
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
