package constpool

import (
    "bufio"
    "container/list"
    "encoding/binary"

    "tklibs/script/assembly"
)

type String struct {
    arrayPool []interface{}
    pool      list.List
}

var _ assembly.ConstPool = &String{}

func (cp *String) Get(index int) interface{} {
    if int(index) < len(cp.arrayPool) {
        return cp.arrayPool[index]
    }

    i := int(index) - len(cp.arrayPool)

    for it := cp.pool.Front(); it != nil; it = it.Next() {
        if i == index {
            return it.Value
        }
        i++
    }

    return nil
}

func (cp *String) Insert(value interface{}) int {
    i := len(cp.arrayPool)
    for it := cp.pool.Front(); it != nil; it = it.Next() {
        if it.Value == value {
            return i
        }
        i++
    }
    cp.pool.PushBack(value)
    return i
}

func (cp *String) clear() {
    cp.arrayPool = make([]interface{}, 0)
    cp.pool.Init()
}

func (cp *String) Write(writer *bufio.Writer) {
    binary.Write(writer, binary.LittleEndian, uint32(cp.pool.Len()+len(cp.arrayPool)))

    for it := cp.pool.Front(); it != nil; it = it.Next() {
        writer.WriteString(it.Value.(string))
        writer.WriteByte(0)
    }
}

func (cp *String) CopyFrom(constPool assembly.ConstPool) {
    src := constPool.(*String)
    cp.arrayPool = src.arrayPool
}

func (cp *String) Read(reader *bufio.Reader) {
    cp.clear()

    var l uint32
    binary.Read(reader, binary.LittleEndian, &l)

    cp.arrayPool = make([]interface{}, l)

    for i := 0; i < int(l); i++ {
        value, err := reader.ReadBytes(0)
        if err != nil {
            return
        }

        cp.arrayPool[i] = string(value[:len(value)-1])
    }
}
