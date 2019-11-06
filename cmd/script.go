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
    useBinary := flag.Bool("b", false, "read binary code from a file")

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

    if !*useBinary {
        scriptCompiler := &struct {
            *compiler.Component
        }{}
        scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)
        for _, file := range files {
            scriptCompiler.AddFile(file)
        }
        asm, tokenList, _ := scriptCompiler.Compile()
        asm.(script.Assembly).Save(bufio.NewWriter(buffer))

        if *showToken {
            for it := tokenList.Front(); it != nil; it = it.Next() {
                fmt.Println(it.Value)
            }
        }
    } else {
        buf, err := ioutil.ReadFile(files[0])
        if err != nil {
            fmt.Println("io error: ", err)
            return
        }
        buffer = bytes.NewBuffer(buf)
    }

    if len(*write) > 0 {
        ioutil.WriteFile(*write, buffer.Bytes(), os.ModePerm)
    }

    executeAsm.Load(bufio.NewReader(buffer))

    if *showAsm {
        for _, inst := range executeAsm.GetFunctions() {
            _func := inst.(runtime.Function)

            fmt.Printf("%v \"%v\", local: %v, args: %v, refs: %v, members: %v\n",
                "Func:",
                _func.GetName(),
                fmt.Sprint(len(_func.GetLocalVars())),
                fmt.Sprint(len(_func.GetArguments())),
                fmt.Sprint(len(_func.GetRefVars())),
                fmt.Sprint(len(_func.GetMembers())))


            //for i, name := range _func.GetMembers() {
            //    fmt.Printf("[%v]Member: %v\n", i, name)
            //}

            fmt.Println(_func.DumpString())
        }
    }

    if !*noExecute {
        println("Begin to execute ...")
        scriptContext := &struct {
            *context.Component
        }{}
        scriptContext.Component = context.NewScriptContext(scriptContext, executeAsm, 4096)

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
