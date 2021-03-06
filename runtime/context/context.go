package context

import (
	"fmt"
	"reflect"
	"unicode"

	"tklibs/script/library/logger"
	"tklibs/script/runtime/util"

	"tklibs/script"
	"tklibs/script/library"
	"tklibs/script/library/basic"
	"tklibs/script/library/debug"
	"tklibs/script/library/io"
	"tklibs/script/library/json"
	"tklibs/script/library/math"
	"tklibs/script/library/semver"
	"tklibs/script/runtime"
	"tklibs/script/runtime/interpreter"
	"tklibs/script/runtime/native"
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

	stringPool util.StringPool

	arrayPrototype *prototype.Array
	mapPrototype   *prototype.Map

	frames          []interface{}
	scopes          []interface{}
	rootRuntimeType runtime.TypeInfo
	assembly        interface{}
	registers       []script.Value
	registerList    [][]script.Value

	globalFields map[string]*script.Value

	initialized     bool
	isProtectObject bool
}

func (impl *Component) GetStringPool() util.StringPool {
	return impl.stringPool
}

func (impl *Component) GetArrayPrototype() interface{} {
	return impl.arrayPrototype
}

func (impl *Component) GetMapPrototype() interface{} {
	return impl.mapPrototype
}

func (impl *Component) IsProtectObject() bool {
	return impl.isProtectObject
}

func (impl *Component) SetProtectObject(value bool) {
	impl.isProtectObject = value
}

var _ runtime.ScriptContext = &Component{}
var _ script.Function = &Component{}
var _ script.MemoryBlock = &Component{}

func (impl *Component) Size() int {
	return 0
}

func (impl *Component) Visit(memoryMap map[interface{}]int, f func(block script.MemoryBlock)) {
	if _, ok := memoryMap[impl]; ok {
		return
	}

	memoryMap[impl] = impl.MemorySize()
	f(impl)

	impl.Component.Visit(memoryMap, f)

	for _, value := range impl.globalFields {
		if ms, ok := value.Get().(script.MemoryBlock); ok {
			ms.Visit(memoryMap, f)
		}
	}
}

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

func (impl *Component) GetRefByName(name string, valuePtr **script.Value) {
	for i := len(impl.scopes) - 1; i >= 0; i-- {
		scope := impl.scopes[i].(runtime.Scope)
		f := scope.GetFunction().(script.Function)
		_func := f.GetScriptRuntimeFunction()

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

	if val, ok := impl.globalFields[name]; ok {
		*valuePtr = val
	} else {
		fieldIndex := -1
		obj := interface{}(impl.functionComponent).(runtime.Object)
		for fieldIndex < 0 && obj != nil {
			fieldIndex = obj.GetRuntimeTypeInfo().(runtime.TypeInfo).GetFieldIndexByName(name)
			if fieldIndex > -1 {
				break
			}
			_prototype := obj.GetPrototype().GetInterface()
			if _prototype != nil {
				obj = _prototype.(runtime.Object)
			} else {
				obj = nil
			}
		}

		if fieldIndex < 0 {
			val := new(script.Value)
			*val = script.NullValue
			*valuePtr = val
			impl.globalFields[impl.stringPool.Insert(name)] = val
		} else {
			*valuePtr = obj.GetByIndex(fieldIndex)
		}
	}
}

func (impl *Component) ScriptSet(fieldName string, value script.Value) {
	impl.globalFields[impl.stringPool.Insert(fieldName)] = &value
}

func (impl *Component) ScriptGet(fieldName string) script.Value {
	if ret, ok := impl.globalFields[fieldName]; ok {
		return *ret
	}

	return impl.functionComponent.ScriptGet(fieldName)
}

func (impl *Component) NewScriptObject(fieldCap int) interface{} {
	n := &struct {
		*object.Component
	}{}

	n.Component = object.NewScriptObject(n, impl, fieldCap)
	return n
}

func (impl *Component) NewScriptArray(sizeCap int) interface{} {
	return array.NewScriptArray(impl.GetOwner(), sizeCap)
}

func (impl *Component) NewScriptMap(sizeCap int) interface{} {
	return _map.NewScriptMap(impl.GetOwner(), sizeCap)
}

func (impl *Component) GetRootRuntimeType() runtime.TypeInfo {
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

func (impl *Component) GetStackFrames() []interface{} {
	return impl.frames
}

func (impl *Component) GetAssembly() interface{} {
	return impl.assembly
}

func (impl *Component) Run() interface{} {
	if !impl.initialized {
		impl.functionComponent.Init()
		impl.initialized = true
	}

	return impl.functionComponent.Invoke(impl.GetOwner(), impl.GetOwner())
}

func (impl *Component) RunWithAssembly(assembly interface{}) interface{} {
	if assembly != nil {
		if _, ok := assembly.(script.Assembly); !ok {
			panic(fmt.Errorf("incorrect assembly type with assembly param"))
		}

		impl.assembly = assembly
	}

	parent := impl.functionComponent
	impl.functionComponent = function.NewScriptFunction(impl.GetOwner(), impl.assembly.(script.Assembly).GetEntry(), impl)
	impl.functionComponent.SetPrototype(script.InterfaceToValue(parent))
	impl.initialized = false
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
			_func := function.NewNativeFunction(*funcDelegate, impl)
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
			_func := function.NewNativeFunction(*funcDelegate, impl)
			obj.(script.Object).ScriptSet(string(name), script.InterfaceToValue(_func))
		}

		impl.functionComponent.ScriptSet(libraryName, script.InterfaceToValue(obj))
	}
}

func (impl *Component) MemoryInfo() runtime.ScriptMemroryInfo {
	memSize := 0

	objCount := 0
	mapCount := 0
	arrayCount := 0
	funcCount := 0

	visit := func(mb script.MemoryBlock) {
		memSize += mb.MemorySize()

		switch mb.(type) {
		case script.Map:
			mapCount++
		case script.Array:
			arrayCount++
		case script.Function:
			funcCount++
		case script.Int:
		case script.Int64:
		case script.Float:
		case script.Float64:
		case script.String:
		case script.Bool:
		case script.Object:
			objCount++
		}
	}

	impl.Visit(make(map[interface{}]int), visit)

	return runtime.ScriptMemroryInfo{
		"mapCount":    mapCount,
		"arrayCount":  arrayCount,
		"funcCount":   funcCount,
		"objectCount": objCount,
		"memorySize":  memSize,
	}
}

func NewScriptContext(owner, asm interface{}, stackSize int) *Component {
	context := &Component{
		ComponentType: script.MakeComponentType(owner),
		assembly:      asm,
		registers:     make([]script.Value, stackSize),
		globalFields:  make(map[string]*script.Value),
		stringPool:    util.NewStringPool(),
	}

	context.registerList = make([][]script.Value, 0, 1)
	context.registerList = append(context.registerList, context.registers)

	runtimeType := typeinfo.NewTypeComponent(context)
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

	context.Object.InitPrototype()
	context.Function.InitPrototype()

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
	context.RegisterLibrary(logger.NewLibrary())
	context.RegisterLibrary(semver.NewLibrary())

	return context
}
