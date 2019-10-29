package benchmark

import (
    "bufio"
    "bytes"
    "testing"

    "tklibs/script"
    "tklibs/script/assembly/assembly"
    "tklibs/script/compiler/compiler"
    "tklibs/script/runtime/context"
)

const (
    DefaultBufSize = 4096
)

func Benchmark_Call(b *testing.B) {
    buffer := bytes.NewBuffer(make([]byte, 0, DefaultBufSize))

    executeAsm := &struct {
        *assembly.Component
    }{}
    executeAsm.Component = assembly.NewScriptAssembly(executeAsm)

    scriptCompiler := &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

    scriptCompiler.AddSource(`
global fib = function(x) {
    if (x <= 0 ) {
        return 0
    }

    if (x == 1 ) {
        return 1
    }

    return fib(x - 1) + fib(x - 2)
}

var result = fib(35)`)

    asm, _, _ := scriptCompiler.Compile()
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))

    executeAsm.Load(bufio.NewReader(buffer))

    scriptContext := &struct {
        *context.Component
    }{}
    scriptContext.Component = context.NewScriptContext(scriptContext, executeAsm, 4096)

    b.ResetTimer()

    for i:=0; i < b.N; i++ {
        scriptContext.Run()
    }
}