package loader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"tklibs/script"
	"tklibs/script/compiler/compiler"
)

//LoadAssemblySourceFromBuffer just load 'source file' from reader mainly for testing
func LoadAssemblySourceFromBuffer(assembly interface{}, reader ...io.Reader) error {
	scriptCompiler := &struct{ *compiler.Component }{}
	scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

	buf := new(bytes.Buffer)
	for _, r := range reader {
		buf.Reset()
		if _, err := buf.ReadFrom(r); err != nil {
			return err
		}
		b := buf.Bytes()
		scriptCompiler.AddSource(*(*string)(unsafe.Pointer(&b)))
	}

	asm, _, err := scriptCompiler.Compile()
	if err != nil {
		return err
	}
	buffer := &bytes.Buffer{}
	asm.(script.Assembly).Save(bufio.NewWriter(buffer))
	assembly.(script.Assembly).Load(bufio.NewReader(buffer))
	return nil
}

func LoadAssembly(assembly interface{}, files ...string) (err error) {
	scriptCompiler := &struct{ *compiler.Component }{}
	scriptCompiler.Component = compiler.NewCompiler(scriptCompiler)

	var filePathGot string
	var fileHandler *os.File
	for _, file := range files {
		// try to find the file
		filePathGot = file
		if !strings.HasSuffix(filePathGot, script.SuffixSourceFile) && !strings.HasSuffix(filePathGot, script.SuffixBinaryCode) {
			filePathGot = file + script.SuffixSourceFile
			if _, err = os.Stat(filePathGot); err != nil {
				filePathGot = file + script.SuffixBinaryCode
				if _, err = os.Stat(filePathGot); err != nil {
					return fmt.Errorf("file:%s not found with suffixs:[.tsc, .tsb]", file)
				}
			}
		} else if _, err = os.Stat(filePathGot); err != nil {
			return fmt.Errorf("file:%s not found", file)
		}

		if strings.HasSuffix(filePathGot, script.SuffixSourceFile) {
			scriptCompiler.AddFile(filePathGot)
		} else {
			if fileHandler, err = os.Open(filePathGot); err == nil {
				assembly.(script.Assembly).Load(bufio.NewReader(fileHandler))
			}
			// we just need only one tsb file
			return
		}
	}

	// only for tsc files
	asm, _, err := scriptCompiler.Compile()
	if err != nil {
		return err
	}
	buffer := &bytes.Buffer{}
	asm.(script.Assembly).Save(bufio.NewWriter(buffer))
	assembly.(script.Assembly).Load(bufio.NewReader(buffer))

	return
}
