package constpool

import (
    "bufio"
    "container/list"
    "encoding/binary"

    "tklibs/script"
    "tklibs/script/assembly"
)

type Int struct {
    pool list.List
}

func (cp *Int) CopyFrom(constPool assembly.ConstPool) {

}

var _ assembly.ConstPool = &Int{}

func (cp *Int) Get(index int) interface{} {
    i := 0
    for it := cp.pool.Front(); it != nil; it = it.Next() {
        if i == index {
            return it.Value
        }
        i++
    }
    return nil
}

func (cp *Int) Insert(value interface{}) int {
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

func (cp *Int) Write(writer *bufio.Writer) {
    binary.Write(writer, binary.LittleEndian, uint32(cp.pool.Len()))
    for it := cp.pool.Front(); it != nil; it = it.Next() {
        binary.Write(writer, binary.LittleEndian, int64(it.Value.(script.Int64)))
    }
}

func (cp *Int) Read(reader *bufio.Reader) {
    var len uint32
    binary.Read(reader, binary.LittleEndian, &len)

    var value int64

    for len > 0 {
        len--
        binary.Read(reader, binary.LittleEndian, &value)
        cp.pool.PushBack(script.Int64(value))
    }
}
