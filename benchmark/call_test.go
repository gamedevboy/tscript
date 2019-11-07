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

func Test_Reload(t *testing.T) {
    source1 := `
func test() {
    println("1")
}
`

    source2 := `
func test() {
    println(toInt("2"))
}
`

    buffer := bytes.NewBuffer(make([]byte, 0, DefaultBufSize))

    executeAsm := &struct {
        *assembly.Component
    }{}
    executeAsm.Component = assembly.NewScriptAssembly(executeAsm)

    scriptCompiler := &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

    scriptCompiler.AddSource(source1)

    asm, _, _ := scriptCompiler.Compile()
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))

    executeAsm.Load(bufio.NewReader(buffer))

    scriptContext := &struct {
        *context.Component
    }{}
    scriptContext.Component = context.NewScriptContext(scriptContext, executeAsm, 4096)
    scriptContext.Run()

    f := scriptContext.ScriptGet("test").GetFunction()
    f.Invoke(nil);

    // reload
    buffer = bytes.NewBuffer(make([]byte, 0, DefaultBufSize))

    scriptCompiler = &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

    scriptCompiler.AddSource(source2)

    asm, _, _ = scriptCompiler.Compile()
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))

    executeAsm2 := &struct {
        *assembly.Component
    }{}
    executeAsm2.Component = assembly.NewScriptAssembly(executeAsm2)
    executeAsm2.Load(bufio.NewReader(buffer))

    scriptContext.ReloadAssembly(executeAsm2)

    f.Invoke(nil)
    //scriptContext.Reload(executeAsm)
}



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