package assembly

import (
    "bufio"
)

type ConstPool interface {
    Insert(interface{}) int
    Get(int) interface{}
    Write(writer *bufio.Writer)
    Read(reader *bufio.Reader)
    CopyFrom(constPool ConstPool)
}
