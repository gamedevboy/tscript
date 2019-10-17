package debug

import (
    "fmt"
    runtime2 "runtime"

    "tklibs/script"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/runtime/stack"
)

type library struct {
    context    interface{}
    Breakpoint native.FunctionType
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
    l.Breakpoint = func(this interface{}, args ...interface{}) interface{} {
        ctx := l.context.(runtime.ScriptContext)
        frame := ctx.GetCurrentFrame().(stack.Frame)
        rf := frame.GetFunction().(runtime.Function)

        i := ctx.(runtime.ScriptInterpreter)

        pc := i.GetPC()
        line := -1

        debugInfo := rf.GetDebugInfoList()
        for i, d := range debugInfo {
            if d.PC > uint32(pc) {
                if i > 0 {
                    line = int(debugInfo[i-1].Line)
                } else {
                    line = int(d.Line)
                }
                break
            }
        }

        if line == -1 {
            line = int(debugInfo[len(debugInfo)-1].Line)
        }

        println()
        //fmt.Printf("SCRIPT BREAKPOINT %q, PC: %v, Line: %v", rf.GetSourceName(), pc, line)
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
