package benchmark

import (
    "bufio"
    "bytes"
    "fmt"
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
ver = 1
func test() {
    println("1")
}
println(ver)
`

    source2 := `
ver = 2
func test() {
    println(toInt("2"))
}
println(ver)
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

    //f := scriptContext.ScriptGet("test").GetFunction()

    version := scriptContext.ScriptGet("ver")
    fmt.Println("source1:", version)
    //f.Invoke(nil);

    // do RunWithAssembly here
    {
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
        scriptContext.RunWithAssembly(executeAsm2)
    }

    //scriptContext.ReloadAssembly(executeAsm2)
    //
    //f.Invoke(nil)

    fmt.Println("source2:", version)
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

var result = fib(35)
println(result)
`)

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