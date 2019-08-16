package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "time"

    "tklibs/script"
    "tklibs/script/assembly/assembly"
    "tklibs/script/compiler/compiler"
    "tklibs/script/runtime"
    "tklibs/script/runtime/context"
)

const (
    DefaultBufSize = 4096
)

func main() {
    showAsm := flag.Bool("s", false, "show asm code")
    showToken := flag.Bool("token", false, "show token")
    noExecute := flag.Bool("n", false, "no execute")
    write := flag.String("w", "", "write binary code to a file")

    flag.Parse()
    files := flag.Args()

    if len(files) < 1 {
        flag.Usage()
        return
    }

    buffer := bytes.NewBuffer(make([]byte, 0, DefaultBufSize))

    executeAsm := &struct {
        *assembly.Component
    }{}
    executeAsm.Component = assembly.NewScriptAssembly(executeAsm)

    scriptCompiler := &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)
    for _, file := range files {
        scriptCompiler.AddFile(file)
    }
    asm, tokenList, _ := scriptCompiler.Compile()
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))

    if len(*write) > 0 {
        ioutil.WriteFile(*write, buffer.Bytes(), os.ModePerm)
    }

    if *showToken {
        for it := tokenList.Front(); it != nil; it = it.Next() {
            fmt.Println(it.Value)
        }
    }

    executeAsm.Load(bufio.NewReader(buffer))

    if *showAsm {
        for _, inst := range executeAsm.GetFunctions() {
            _func := inst.(runtime.Function)

            fmt.Printf("%v \"%v\", local: %v, args: %v, refs: %v, members: %v",
                "Func:",
                _func.GetName(),
                fmt.Sprint(len(_func.GetLocalVars())),
                fmt.Sprint(len(_func.GetArguments())),
                fmt.Sprint(len(_func.GetRefVars())),
                fmt.Sprint(len(_func.GetMembers())))

            println()

            print(_func.DumpString())
            println()
        }
    }

    if !*noExecute {
        println("Begin to execute ...")
        scriptContext := &struct {
            *context.Component
        }{}
        scriptContext.Component = context.NewScriptContext(scriptContext, executeAsm)

        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("Error: %v", fmt.Sprint(err))
            }
        }()

        startTime := time.Now()
        scriptContext.Run()
        fmt.Printf("Elasped time: %v ms", fmt.Sprint(time.Since(startTime).Nanoseconds()/1000000))
    }
}
