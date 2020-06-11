package debug

import (
    "fmt"
    runtime2 "runtime"

    "tklibs/script"
	"tklibs/script/library/logger"
	"tklibs/script/runtime"
    "tklibs/script/runtime/native"
    "tklibs/script/runtime/runtime_t"
    "tklibs/script/runtime/stack"
)

type CallInfo struct {
	FilePath string
	Line     int
	FuncName string
}

func GetCallInfo(sc runtime.ScriptContext) *CallInfo {
	frame := sc.GetCurrentFrame().(stack.Frame)
	rf := frame.GetFunction().(runtime_t.Function)
	pc := sc.(runtime.ScriptInterpreter).GetPC()
	debugInfo := rf.GetDebugInfoList()
	debugInfoLen := len(debugInfo)
	sourceIndex := -1
	line := -1
	for i, d := range debugInfo {
		if d.PC > uint32(pc) {
			if i > 0 {
				line = int(debugInfo[i-1].Line)
				sourceIndex = int(debugInfo[i-1].SourceIndex)
			} else {
				line = int(d.Line)
				sourceIndex = int(d.SourceIndex)
			}
			break
		}
	}

	if line == -1 {
		line = int(debugInfo[debugInfoLen-1].Line)
	}

	if sourceIndex == -1 {
		sourceIndex = int(debugInfo[debugInfoLen-1].SourceIndex)
	}
	return &CallInfo{FilePath: rf.GetSourceNames()[sourceIndex], Line: line, FuncName: rf.GetName()}
}

type library struct {
	context    interface{}
	Breakpoint native.FunctionType
	Log        native.FunctionType
}

func (*library) GetName() string {
	return "debug"
}

func (l *library) SetScriptContext(context interface{}) {
	l.context = context
}

func NewLibrary() *library {
	ret := &library{}
	ret.init()
	return ret
}

func (l *library) init() {
	l.Log = func(this interface{}, args ...interface{}) interface{} {
		logger.ScriptLogger().Debug(args...)
		return script.Null
	}

	l.Breakpoint = func(this interface{}, args ...interface{}) interface{} {
		ctx := l.context.(runtime.ScriptContext)
		i := ctx.(runtime.ScriptInterpreter)

		_ = GetCallInfo(ctx)

		println()

		for i, v := range i.GetCurrentRegisters() {
			switch t := v.Get().(type) {
			case script.Int:
				fmt.Printf("[%v] \t%v", i, t)
			case script.Float:
				fmt.Printf("[%v] \t%v", i, t)
			case script.String:
				fmt.Printf("[%v] \t%v", i, t)
			case script.Bool:
				fmt.Printf("[%v] \t%v", i, t)
			case script.Object:
				fmt.Printf("[%v] \tObject:%v", i, t.GetScriptTypeId())
			default:
				fmt.Printf("[%v] \t%v", i, t)
			}

			println()
		}

		runtime2.Breakpoint()

		return script.Null
	}
}
