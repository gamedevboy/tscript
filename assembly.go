package script

import (
    "bufio"

    "tklibs/script/assembly"
)

type Assembly interface {
    GetStringConstPool() assembly.ConstPool
    GetIntConstPool() assembly.ConstPool
    GetFloatConstPool() assembly.ConstPool
    GetEntry() interface{}
    Load(reader *bufio.Reader)
    Save(writer *bufio.Writer)
    GetFunctionByMetaIndex(Int) interface{}
    GetFunctions() []interface{}
    Reload(assembly Assembly) error
	Update()
}
