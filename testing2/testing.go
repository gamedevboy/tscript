package testing2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"tklibs/script"
	"tklibs/script/assembly/assembly"
	"tklibs/script/assembly/loader"
	"tklibs/script/compiler/compiler"
	"tklibs/script/library/logger"
	"tklibs/script/runtime/context"
)

func MustCompile(sources ...string) *bytes.Buffer{
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	scriptCompiler := &struct {
		*compiler.Component
	}{}
	scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)
	for _, file := range sources {
		scriptCompiler.AddSource(file)
	}
	asm, _, err := scriptCompiler.Compile()
	if err != nil {
		panic(err)
	}
	asm.(script.Assembly).Save(bufio.NewWriter(buffer))
	return buffer
}

func MustCompileToTemp(sources ...string) string{
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	scriptCompiler := &struct {
		*compiler.Component
	}{}
	scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)
	for _, s := range sources {
		scriptCompiler.AddSource(s)
	}
	asm, _, err := scriptCompiler.Compile()
	if err != nil {
		panic(err)
	}
	asm.(script.Assembly).Save(bufio.NewWriter(buffer))

	tsb := fmt.Sprintf(path.Join(os.TempDir(), fmt.Sprintf("script_must_compile_to_temp_%d.tsb",time.Now().UnixNano())))
	if err := ioutil.WriteFile(tsb,buffer.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
	logger.ScriptLogger().Debugf("tsb compile to: %s",tsb)
	return tsb
}

func MustInitWithSource(sources ...string) (*context.Component,*assembly.Component) {
	reader := make([]io.Reader, 0)
	for _, s := range sources {
		reader = append(reader, strings.NewReader(s))
	}

	scriptAssembly := &struct{ *assembly.Component }{}
	scriptAssembly.Component = assembly.NewScriptAssembly(scriptAssembly)
	if err := loader.LoadAssemblySourceFromBuffer(scriptAssembly, reader...); err != nil {
		panic(err)
	}
	scriptContext := &struct{ *context.Component }{}
	scriptContext.Component = context.NewScriptContext(scriptContext, scriptAssembly, 64)
	return scriptContext.Component,scriptAssembly.Component
}

func MustInitWithFile(filePath ...string) (*context.Component,*assembly.Component) {
	scriptAssembly := &struct{ *assembly.Component }{}
	scriptAssembly.Component = assembly.NewScriptAssembly(scriptAssembly)
	if err := loader.LoadAssembly(scriptAssembly, filePath...); err != nil {
		panic(err)
	}
	scriptContext := &struct{ *context.Component }{}
	scriptContext.Component = context.NewScriptContext(scriptContext, scriptAssembly, 64)
	scriptContext.Run()
	return scriptContext.Component,scriptAssembly.Component
}
