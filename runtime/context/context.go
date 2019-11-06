package context

import (
    "fmt"
    "reflect"
    "unicode"

    "tklibs/script"
    "tklibs/script/library"
    "tklibs/script/library/basic"
    "tklibs/script/library/debug"
    "tklibs/script/library/io"
    "tklibs/script/library/json"
    "tklibs/script/library/math"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/runtime/interpreter"
    "tklibs/script/runtime/prototype"
    "tklibs/script/runtime/typeinfo"
    "tklibs/script/type/array"
    "tklibs/script/type/function"
    _map "tklibs/script/type/map"

    "tklibs/script/type/object"
)

type functionComponent = function.Component
type interpreterComponent = interpreter.Component

type Component struct {
    script.ComponentType

    *functionComponent
    *interpreterComponent

    *prototype.Bool
    *prototype.Function
    *prototype.Number
    *prototype.Object
    *prototype.String

    arrayPrototype    *prototype.Array
    mapPrototype      *prototype.Map

    frames          []interface{}
    scopes          []interface{}
    rootRuntimeType interface{}
    assembly        interface{}
    registers       []script.Value
    registerList    [][]script.Value
}

func (impl *Component) GetArrayPrototype() interface{} {
    return impl.arrayPrototype
}

func (impl *Component) GetMapPrototype() interface{} {
    return impl.mapPrototype
}

var _ runtime.ScriptContext = &Component{}

func (impl *Component) PushRegisters(regStart script.Int, length int) []script.Value {
    if len(impl.registers[regStart:]) < length {
        impl.registers = append(impl.registers, make([]script.Value, length-len(impl.registers[regStart:]))...)
    }

    impl.registerList = append(impl.registerList, impl.registers)
    r := impl.registers
    impl.registers = impl.registers[regStart:]
    return r
}

func (impl *Component) PopRegisters() {
    impl.registers = impl.registerList[len(impl.registerList)-1]
    impl.registerList = impl.registerList[:len(impl.registerList)-1]
}

func (impl *Component) GetRegisters() []script.Value {
    return impl.registers
}

func (impl *Component) GetRefByName(name string, valuePtr **script.Value)  {
    for i := len(impl.scopes) - 1; i >= 0; i-- {
        scope := impl.scopes[i].(runtime.Scope)
        f := scope.GetFunction().(script.Function)
        _func := f.GetRuntimeFunction().(runtime.Function)

        for i, n := range _func.GetLocalVars() {
            if name == n {
                scope.AddToRefList(&scope.GetLocalVarList()[i], valuePtr)
                return
            }
        }

        for i, n := range _func.GetArguments() {
            if name == n {
                scope.AddToRefList(&scope.GetArgList()[i], valuePtr)
                return
            }
        }
    }

    fieldIndex := impl.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldIndexByName(name)

    if fieldIndex < 0 {
        impl.ScriptSet(name, script.NullValue)
        fieldIndex = impl.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldIndexByName(name)
    }

    *valuePtr = impl.GetByIndex(fieldIndex)

    return
}

func (impl *Component) NewScriptObject(fieldCap int) interface{} {
    n := &struct {
        *object.Component
    }{}

    n.Component = object.NewScriptObject(n, impl, fieldCap)
    return n
}

func (impl *Component) NewScriptArray(sizeCap int) interface{} {
    a := &struct {
        *array.Component
    }{}
    a.Component = array.NewScriptArray(a, impl.GetOwner(), sizeCap)
    return a
}

func (impl *Component) NewScriptMap(sizeCap int) interface{} {
    m := &struct {
        *_map.Component
    }{}
    m.Component = _map.NewScriptMap(m, impl.GetOwner(), sizeCap)
    return m
}

func (impl *Component) GetRootRuntimeType() interface{} {
    return impl.rootRuntimeType
}

func (impl *Component) PushFrame(frame interface{}) {
    impl.frames = append(impl.frames, frame)
}

func (impl *Component) PopFrame() interface{} {
    ret := impl.GetCurrentFrame()
    impl.frames = impl.frames[:len(impl.frames)-1]
    return ret
}

func (impl *Component) PushScope(value interface{}) {
    impl.scopes = append(impl.scopes, value)
}

func (impl *Component) PopScope() interface{} {
    back := impl.scopes[len(impl.scopes)-1]
    back.(runtime.Scope).KeepRefs()
    impl.scopes = impl.scopes[:len(impl.scopes)-1]
    return back
}

func (impl *Component) GetCurrentFrame() interface{} {
    return impl.frames[len(impl.frames)-1]
}

func (impl *Component) GetAssembly() interface{} {
    return impl.assembly
}

func (impl *Component) Run() interface{} {
    ret := impl.functionComponent.Invoke(impl.GetOwner())
    impl.registers[0].SetNull()
    return ret;
}

func (impl *Component) Reload(assembly interface{}) interface{} {
    if assembly != nil {
        if _, ok := assembly.(script.Assembly); !ok {
            panic(fmt.Errorf("incorrect assembly type with assembly param"))
        }

        impl.assembly = assembly
    }

    parent := impl.functionComponent
    impl.functionComponent = function.NewScriptFunction(impl.GetOwner(), impl.assembly.(script.Assembly).GetEntry(), impl)
    impl.functionComponent.SetPrototype(script.InterfaceToValue(parent))
    return impl.Run()
}

func (impl *Component) RegisterLibrary(library library.RuntimeLibrary) {
    library.SetScriptContext(impl.GetOwner())
    libraryName := library.GetName()
    switch libraryName {
    case "":
        valueOfBasicLibrary := reflect.ValueOf(library).Elem()
        libraryType := valueOfBasicLibrary.Type()
        funcDelegate := new(native.FunctionType)
        valueOfDelegate := reflect.ValueOf(funcDelegate).Elem()

        for i := 0; i < valueOfBasicLibrary.NumField(); i++ {
            if libraryType.Field(i).Type.Kind() != reflect.Func {
                continue
            }
            name := []rune(libraryType.Field(i).Name)
            name[0] = unicode.ToLower(name[0])
            valueOfDelegate.Set(valueOfBasicLibrary.Field(i))
            _func := native.NewNativeFunction(*funcDelegate, impl)
            impl.functionComponent.ScriptSet(string(name), script.InterfaceToValue(_func))
        }
    default:
        obj := impl.NewScriptObject(0)

        valueOfBasicLibrary := reflect.ValueOf(library).Elem()
        libraryType := valueOfBasicLibrary.Type()
        funcDelegate := new(native.FunctionType)
        valueOfDelegate := reflect.ValueOf(funcDelegate).Elem()

        for i := 0; i < valueOfBasicLibrary.NumField(); i++ {
            if libraryType.Field(i).Type.Kind() != reflect.Func {
                continue
            }
            name := []rune(libraryType.Field(i).Name)
            name[0] = unicode.ToLower(name[0])
            valueOfDelegate.Set(valueOfBasicLibrary.Field(i))
            _func := native.NewNativeFunction(*funcDelegate, impl)
            obj.(script.Object).ScriptSet(string(name), script.InterfaceToValue(_func))
        }

        impl.functionComponent.ScriptSet(libraryName, script.InterfaceToValue(obj))
    }
}

func NewScriptContext(owner, asm interface{}, stackSize int) *Component {
    context := &Component{
        ComponentType: script.MakeComponentType(owner),
        assembly:      asm,
        registers:     make([]script.Value, stackSize),
    }

    context.registerList = make([][]script.Value, 0, 1)
    context.registerList = append(context.registerList, context.registers)

    runtimeType := &struct {
        *typeinfo.Component
    }{}
    runtimeType.Component = typeinfo.NewTypeComponent(runtimeType)
    context.rootRuntimeType = runtimeType

    context.Object = prototype.NewObjectPrototype(context)
    context.Function = prototype.NewFunctionPrototype(context)
    context.Bool = prototype.NewBoolPrototype(context)
    context.Number = prototype.NewNumberPrototype(context)
    context.String = prototype.NewStringPrototype(context)
    context.arrayPrototype = prototype.NewArrayPrototype(context)
    context.mapPrototype = prototype.NewMapPrototype(context)

    context.functionComponent = function.NewScriptFunction(owner, asm.(script.Assembly).GetEntry(), context)
    context.interpreterComponent = interpreter.NewScriptInterpreter(owner, context)

    context.Object.Init()
    context.Function.Init()

    context.ScriptSet("Object", script.InterfaceToValue(context.GetObjectPrototype()))
    context.ScriptSet("Array", script.InterfaceToValue(context.GetArrayPrototype()))
    context.ScriptSet("String", script.InterfaceToValue(context.GetStringPrototype()))
    context.ScriptSet("Map", script.InterfaceToValue(context.GetMapPrototype()))
    context.ScriptSet("Number", script.InterfaceToValue(context.GetNumberPrototype()))

    context.ScriptSet("$G", script.InterfaceToValue(context))

    context.RegisterLibrary(basic.NewLibrary())
    context.RegisterLibrary(json.NewLibrary())
    context.RegisterLibrary(io.NewLibrary())
    context.RegisterLibrary(math.NewLibrary())
    context.RegisterLibrary(debug.NewLibrary())

    return context
}
