package script

import "fmt"

type Error interface {
    GetLine() int
    GetFilePath() string
}

type scriptError struct {
    message  string
    line     int
    filePath string
}

func (e *scriptError) GetLine() int {
    return e.line
}

func (e *scriptError) GetFilePath() string {
    return e.filePath
}

func (e *scriptError) Error() string {
    return e.message
}

var _ error = &scriptError{}
var _ Error = &scriptError{}

func MakeError(filePath string, line int, format string, args ...interface{}) error {
    return &scriptError{
        line:     line,
        filePath: filePath,
        message:  fmt.Sprintf(format, args ...),
    }
}
