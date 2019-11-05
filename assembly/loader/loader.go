package loader

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "unsafe"

    "tklibs/script"
    "tklibs/script/compiler/compiler"
)

//LoadAssemblySourceFromBuffer just load 'source file' from reader mainly for testing
func LoadAssemblySourceFromBuffer(assembly interface{},reader ...io.Reader) error {
    scriptCompiler := &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

    buf := new(bytes.Buffer)
    for _, r := range reader {
        buf.Reset()
        if _, err :=buf.ReadFrom(r) ; err != nil {
            return err
        }
        b := buf.Bytes()
        scriptCompiler.AddSource(*(*string)(unsafe.Pointer(&b)))
    }

    asm, _, _ := scriptCompiler.Compile()
    buffer := &bytes.Buffer{}
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))
    assembly.(script.Assembly).Load(bufio.NewReader(buffer))
    return nil
}


func LoadAssembly(assembly interface{}, files ...string) error {
    buffer := &bytes.Buffer{}

    scriptCompiler := &struct {
        *compiler.Component
    }{}
    scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

    for _, file := range files {
        // check source file first
        filePath := file + ".tsc"

        _, err := os.Stat(filePath)
        if err == nil {
            scriptCompiler.AddFile(filePath)
        } else {
            // check source file type
            filePath = file + ".tsb"
            _, err = os.Stat(filePath)
            if err == nil {
                file, err := os.Open(filePath)

                if err == nil {
                    assembly.(script.Assembly).Load(bufio.NewReader(file))
                    return nil
                }

                return err
            }
        }
    }

    asm, _, _ := scriptCompiler.Compile()
    asm.(script.Assembly).Save(bufio.NewWriter(buffer))

    assembly.(script.Assembly).Load(bufio.NewReader(buffer))

    return nil
}
