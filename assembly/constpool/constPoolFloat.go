package constpool

import (
    "bufio"
    "container/list"
    "encoding/binary"

    "tklibs/script"
)

type Float struct {
    pool list.List
}

func (cp *Float) Get(index int) interface{} {
    i := 0
    for it := cp.pool.Front(); it != nil; it = it.Next() {
        if i == index {
            return it.Value
        }
        i++
    }
    return nil
}

func (cp *Float) Insert(value interface{}) int {
    i := 0
    for it := cp.pool.Front(); it != nil; it = it.Next() {
        if it.Value == value {
            return i
        }
        i++
    }
    cp.pool.PushBack(value)
    return i
}

func (cp *Float) Write(writer *bufio.Writer) {
    binary.Write(writer, binary.LittleEndian, uint32(cp.pool.Len()))

    for it := cp.pool.Front(); it != nil; it = it.Next() {
        binary.Write(writer, binary.LittleEndian, it.Value)
    }
}

func (cp *Float) Read(reader *bufio.Reader) {
    var len uint32

    binary.Read(reader, binary.LittleEndian, &len)

    var value script.Float64

    for len > 0 {
        len--
        binary.Read(reader, binary.LittleEndian, &value)
        cp.pool.PushBack(value)
    }
}
