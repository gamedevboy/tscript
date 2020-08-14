package assembly

import (
    "bufio"
    "encoding/binary"
    "fmt"

    "tklibs/script"
    "tklibs/script/assembly"
    "tklibs/script/assembly/constpool"
    "tklibs/script/runtime/runtime_t"
)

type Component struct {
    script.ComponentType
    intConstPool    assembly.ConstPool
    floatConstPool  assembly.ConstPool
    stringConstPool assembly.ConstPool

    functions []interface{}
}

func (impl *Component) Update() {
    for _, f := range impl.functions {
        f.(runtime_t.Function).Update()
    }
}

var _ script.Assembly = &Component{}

func (impl *Component) GetFunctionByMetaIndex(index script.Int) interface{} {
    return impl.functions[index]
}

func (impl *Component) GetFunctions() []interface{} {
    return impl.functions
}

type assemblyFileHeader struct {
    magic   uint32
    version uint32
}

func (impl *Component) GetStringConstPool() assembly.ConstPool {
    return impl.stringConstPool
}

func (impl *Component) GetIntConstPool() assembly.ConstPool {
    return impl.intConstPool
}

func (impl *Component) GetFloatConstPool() assembly.ConstPool {
    return impl.floatConstPool
}

func (impl *Component) GetEntry() interface{} {
    return impl.functions[0]
}

func (impl *Component) Load(reader *bufio.Reader) {
    fh := &assemblyFileHeader{}
    binary.Read(reader, binary.LittleEndian, &fh.magic)
    binary.Read(reader, binary.LittleEndian, &fh.version)

    loadV1(fh, impl, reader)
}

func (impl *Component) Reload(assembly script.Assembly) error {
    if len(impl.GetFunctions()) != len(assembly.GetFunctions()) {
        return fmt.Errorf("Can't reload assembly due to mismatch function count ")
    }

    impl.GetStringConstPool().CopyFrom(assembly.GetStringConstPool())
    impl.GetIntConstPool().CopyFrom(assembly.GetIntConstPool())
    impl.GetFloatConstPool().CopyFrom(assembly.GetFloatConstPool())

    for i, f := range assembly.GetFunctions() {
        dest := impl.functions[i].(runtime_t.Function)
        src := f.(runtime_t.Function)

        if dest.GetName() != src.GetName() {
            panic("")
        }

        dest.CopyFrom(src)
    }

    return nil
}

func (impl *Component) Save(writer *bufio.Writer) {
    fh := &assemblyFileHeader{}
    binary.Write(writer, binary.LittleEndian, &fh.magic)
    binary.Write(writer, binary.LittleEndian, &fh.version)
    saveV1(&assemblyFileHeader{}, impl, writer)

    writer.Flush()
}

func NewScriptAssemblyWithFunctions(owner interface{}, functions []interface{}) *Component {
    return &Component{
        ComponentType:   script.MakeComponentType(owner),
        functions:       functions,
        stringConstPool: &constpool.String{},
        intConstPool:    &constpool.Int{},
        floatConstPool:  &constpool.Float{},
    }
}

func NewScriptAssembly(owner interface{}) *Component {
    return &Component{
        ComponentType:   script.MakeComponentType(owner),
        stringConstPool: &constpool.String{},
        intConstPool:    &constpool.Int{},
        floatConstPool:  &constpool.Float{},
    }
}
