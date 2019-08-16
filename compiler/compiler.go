package compiler

import "container/list"

type Compiler interface {
    AddFile(fileName string)
    Compile() (interface{}, *list.List, error)
}
