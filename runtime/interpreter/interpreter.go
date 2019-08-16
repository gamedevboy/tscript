package interpreter

import (
    "fmt"
    "math"
    "unsafe"

    "tklibs/script"
    "tklibs/script/instruction"
    "tklibs/script/opcode"
    "tklibs/script/runtime"
    "tklibs/script/runtime/function/native"
    "tklibs/script/runtime/scope"
    "tklibs/script/runtime/stack/frame"
    "tklibs/script/runtime/util"
    "tklibs/script/type/function"
    "tklibs/script/value"
)

type Component struct {
    script.ComponentType
    context          interface{}
    currentPC        int
    currentRegisters []script.Value
}

func (impl *Component) GetCurrentRegisters() []script.Value {
    return impl.currentRegisters
}

func (impl *Component) GetPC() int {
    return impl.currentPC
}

var _ runtime.ScriptInterpreter = &Component{}

func (impl *Component) InvokeNew(function interface{}, args ...interface{}) interface{} {
    this := impl.context.(runtime.ScriptContext).NewScriptObject(0)
    this.(runtime.Object).SetPrototype(script.InterfaceToValue(function))
    sf := function.(script.Function)
    switch _func := sf.GetRuntimeFunction().(type) {
    case runtime.Function:
        if len(args) < len(_func.GetArguments()) {
            panic("") //todo not enough arguments
        }
        context := impl.context.(runtime.ScriptContext)
        context.PushRegisters(0, _func.GetMaxRegisterCount()+len(_func.GetLocalVars())+len(_func.GetArguments())+2)
        defer context.PopRegisters()
        registers := context.GetRegisters()
        registers[1].Set(this)
        for i := range _func.GetArguments() {
            registers[2+i].Set(args[i])
        }
        impl.invoke(function)
        return registers[0].Get()
    case native.Function:
        return _func.NativeCall(this, args...)
    default:
        panic("")
    }
}

func (impl *Component) InvokeFunction(function, this interface{}, args ...interface{}) interface{} {
    sf := function.(script.Function)
    switch _func := sf.GetRuntimeFunction().(type) {
    case runtime.Function:
        if len(args) < len(_func.GetArguments()) {
            panic("") //todo not enough arguments
        }
        context := impl.context.(runtime.ScriptContext)
        context.PushRegisters(0, _func.GetMaxRegisterCount()+len(_func.GetLocalVars())+len(_func.GetArguments())+2)
        defer context.PopRegisters()
        registers := context.GetRegisters()
        registers[1].Set(this)
        for i := range _func.GetArguments() {
            registers[2+i].Set(args[i])
        }
        impl.invoke(function)
        return registers[0].Get()
    case native.Function:
        return _func.NativeCall(this, args...)
    default:
        panic("")
    }
}

//noinspection GoNilness
func (impl *Component) invoke(f interface{}) {
    sf := f.(script.Function)
    _func := sf.GetRuntimeFunction().(runtime.Function)
    context := impl.context.(runtime.ScriptContext)
    registers := context.GetRegisters()

    if _func.IsCaptureThis() {
        registers[1] = sf.GetThis()
    }

    instList := _func.GetInstructionList()
    instCount := len(instList)
    if instCount == 0 {
        registers[0].SetInterface(nil)
        return
    }

    _frame := frame.NewStackFrame(nil, _func)
    context.PushFrame(_frame)
    defer func() { frame.FreeStackFrame(context.PopFrame().(*frame.Component)) }()

    if _func.IsScope() {
        s := scope.NewScope(nil, f, registers[2:], registers[ 2+len(_func.GetArguments()):])

        context.PushScope(s)
        defer func() { scope.FreeScope(context.PopScope().(*scope.Component)) }()
    }

    var vb, vc script.Value
    var pa_, pb_, pc_ *script.Value

    ilStart := uintptr(unsafe.Pointer(&instList[0]))

    pc := 0

    defer func() {
        if err := recover(); err != nil {
            switch e := err.(type) {
            case script.Error:
                panic(e)
            default:
                debugInfo := _func.GetDebugInfoList()
                debugInfoLen := len(debugInfo)

                sourceIndex := -1
                line := -1
                for i, d := range debugInfo {
                    if d.PC > uint32(pc) {
                        if i > 0 {
                            line = int(debugInfo[i-1].Line)
                            sourceIndex = int(debugInfo[i-1].SourceIndex)
                        } else {
                            line = int(d.Line)
                            sourceIndex = int(d.SourceIndex)
                        }
                        break
                    }
                }

                if line == -1 {
                    line = int(debugInfo[debugInfoLen-1].Line)
                }

                if sourceIndex == -1 {
                    sourceIndex = int(debugInfo[debugInfoLen-1].SourceIndex)
                }

                fileName := _func.GetSourceNames()[sourceIndex]
                panic(script.MakeError(fileName, line, "script runtime error: line[%v] in %v: %v", line, fileName, err))
            }
        }
    }()

    for pc < instCount {
        il := (*instruction.Instruction)(unsafe.Pointer(ilStart + uintptr(pc*8)))

        // decode the opcode
        _type := il.Type >> 4
        bcType := il.Type & 15

        if il.A > -1 {
            pa_ = &registers[il.A]
        } else {
            pa_ = sf.GetRefList()[-il.A-1]
        }

        switch bcType {
        case opcode.Register | opcode.Register<<2:
            pc_ = &registers[il.C]
            pb_ = &registers[il.B]
        case opcode.Register | opcode.Integer<<2:
            pc_ = &registers[il.C]
            vb.SetInt(script.Int(il.B))
            pb_ = &vb
        case opcode.Integer | opcode.Register<<2:
            vc.SetInt(script.Int(il.C))
            pc_ = &vc
            pb_ = &registers[il.B]
        case opcode.Register | opcode.Reference<<2:
            pc_ = &registers[il.C]
            pb_ = sf.GetRefList()[il.B]
        case opcode.Reference | opcode.Register<<2:
            pc_ = sf.GetRefList()[il.C]
            pb_ = &registers[il.B]
        case opcode.Reference | opcode.Reference<<2:
            pc_ = sf.GetRefList()[il.C]
            pb_ = sf.GetRefList()[il.B]
        case opcode.Reference | opcode.Integer<<2:
            pc_ = sf.GetRefList()[il.C]
            vb.SetInt(script.Int(il.B))
            pb_ = &vb
        case opcode.Integer | opcode.Reference<<2:
            vc.SetInt(script.Int(il.C))
            pc_ = &vc
            pb_ = sf.GetRefList()[il.B]
        case opcode.Integer | opcode.Integer<<2:
            vc.SetInt(script.Int(il.C))
            vb.SetInt(script.Int(il.B))
            pc_ = &vc
            pb_ = &vb
        case opcode.Register << 2:
            pb_ = &registers[il.B]
        case opcode.Reference << 2:
            pb_ = sf.GetRefList()[il.B]
        case opcode.Integer << 2:
            vb.SetInt(il.GetABx().B)
            pb_ = &vb
        case opcode.None:
            b := il.GetABm().B
            if math.IsNaN(float64(b)) {
                if il.GetABx().B == math.MaxInt32 {
                    vb.SetBool(true)
                } else {
                    vb.SetBool(false)
                }
            } else {
                vb.SetFloat(b)
            }
            pb_ = &vb
        }

        switch _type {
        case opcode.Memory:
            switch il.Code {
            case opcode.Move:
                *pa_ = *pb_
            case opcode.LoadField:
                *pa_ = sf.GetFieldByMemberIndex(pb_.Get(), pc_.GetInt())
            case opcode.StoreField:
                sf.SetFieldByMemberIndex(pa_.Get(), pb_.GetInt(), *pc_)
            case opcode.LoadElement:
                switch pb_.GetType() {
                case script.ValueTypeInterface:
                    switch target := pb_.GetInterface().(type) {
                    case script.Map:
                        *pa_ = target.GetValue(*pc_)
                    case script.Array:
                        *pa_ = target.GetElement(pc_.ToInt())
                    case script.Object:
                        *pa_ = target.ScriptGet(string(util.ToScriptString(pc_.Get())))
                    }
                default:
                    panic("")
                }
            case opcode.StoreElement:
                switch pa_.GetType() {
                case script.ValueTypeInterface:
                    switch target := pa_.GetInterface().(type) {
                    case script.Map:
                        target.SetValue(*pb_, *pc_)
                    case script.Array:
                        target.SetElement(pb_.ToInt(), *pc_)
                    case script.Object:
                        target.ScriptSet(string(pb_.GetInterface().(script.String)), *pc_)
                    }
                default:
                    panic("")
                }
            case opcode.Map:
                pa_.SetInterface(context.NewScriptMap(int(pb_.GetInt())))
            case opcode.Array:
                pa_.SetInterface(context.NewScriptArray(int(pb_.GetInt())))
            }
        case opcode.Const:
            switch il.Code {
            case opcode.Load:
                index := pb_.GetInt()
                _t := index & 3
                index = index >> 2
                switch _t {
                case opcode.ConstInt64:
                    pa_.SetInt64(context.GetAssembly().(script.Assembly).GetIntConstPool().Get(int(index)).(script.Int64))
                case opcode.ConstFloat64:
                    pa_.SetFloat64(context.GetAssembly().(script.Assembly).GetFloatConstPool().Get(int(index)).(script.Float64))
                case opcode.ConstString:
                    pa_.SetInterface(script.String(context.GetAssembly().(script.Assembly).GetStringConstPool().Get(int(index)).(string)))
                }
            case opcode.LoadFunc:
                metaIndex := pb_.GetInt()
                f := &struct {
                    *function.Component
                }{}
                f.Component = function.NewScriptFunction(f, context.GetAssembly().(script.Assembly).GetFunctionByMetaIndex(metaIndex),
                    context)
                pa_.SetInterface(f)
                rf := f.GetRuntimeFunction().(runtime.Function)
                for i := 0; i < len(rf.GetRefVars()); i++ {
                    context.GetRefByName(rf.GetRefVars()[i], &f.GetRefList()[i])
                }
                if rf.IsCaptureThis() {
                    f.SetThis(registers[1])
                }
            case opcode.LoadNil:
                pa_.SetInterface(script.Null)
            }
        case opcode.Math:
            switch il.Code {
            case opcode.Add:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetInt(pb_.GetInt() + pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetFloat(script.Float(pb_.GetInt()) + pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetInt64(script.Int64(pb_.GetInt()) + vc_)
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetInt()) + vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetFloat(pb_.GetFloat() + pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetFloat(pb_.GetFloat() + script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) + vc_)
                        case script.Int64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) + script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetInt64(vb_ + vc_)
                            case script.Float64:
                                pa_.SetFloat64(script.Float64(vb_) + vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetInt64(vb_ + script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(script.Float64(vb_) + script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetFloat64(vb_ + vc_)
                            case script.Int64:
                                pa_.SetFloat64(vb_ + script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetFloat64(vb_ + script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(vb_ + script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("+").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic(fmt.Errorf("Add can not support: %v: %v ", pb_.GetType(), pb_.Get()))
                }
            case opcode.Sub:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetInt(pb_.GetInt() - pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetFloat(script.Float(pb_.GetInt()) - pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetInt64(script.Int64(pb_.GetInt()) - vc_)
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetInt()) - vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetFloat(pb_.GetFloat() - pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetFloat(pb_.GetFloat() - script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) - vc_)
                        case script.Int64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) - script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetInt64(vb_ - vc_)
                            case script.Float64:
                                pa_.SetFloat64(script.Float64(vb_) - vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetInt64(vb_ - script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(script.Float64(vb_) - script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetFloat64(vb_ - vc_)
                            case script.Int64:
                                pa_.SetFloat64(vb_ - script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetFloat64(vb_ - script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(vb_ - script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("-").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Mul:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetInt(pb_.GetInt() * pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetFloat(script.Float(pb_.GetInt()) * pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetInt64(script.Int64(pb_.GetInt()) * vc_)
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetInt()) * vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetFloat(pb_.GetFloat() * pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetFloat(pb_.GetFloat() * script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) * vc_)
                        case script.Int64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) * script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetInt64(vb_ * vc_)
                            case script.Float64:
                                pa_.SetFloat64(script.Float64(vb_) * vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetInt64(vb_ * script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(script.Float64(vb_) * script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetFloat64(vb_ * vc_)
                            case script.Int64:
                                pa_.SetFloat64(vb_ * script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetFloat64(vb_ * script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(vb_ * script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("*").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Div:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetInt(pb_.GetInt() / pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetFloat(script.Float(pb_.GetInt()) / pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetInt64(script.Int64(pb_.GetInt()) / vc_)
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetInt()) / vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetFloat(pb_.GetFloat() / pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetFloat(pb_.GetFloat() / script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) / vc_)
                        case script.Int64:
                            pa_.SetFloat64(script.Float64(pb_.GetFloat()) / script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetInt64(vb_ / vc_)
                            case script.Float64:
                                pa_.SetFloat64(script.Float64(vb_) / vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetInt64(vb_ / script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(script.Float64(vb_) / script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        pa_.SetFloat64(vb_ + pc_.ToFloat64())
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetFloat64(vb_ / vc_)
                            case script.Int64:
                                pa_.SetFloat64(vb_ / script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetFloat64(vb_ / script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetFloat64(vb_ / script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("/").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Rem:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetInt(pb_.GetInt() % pc_.GetInt())

                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetInt64(script.Int64(pb_.GetInt()) % vc_)

                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetInt64(vb_ % vc_)

                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetInt64(vb_ % script.Int64(pc_.GetInt()))

                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("%").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Inc:
                switch pa_.GetType() {
                case script.ValueTypeInt:
                    pa_.SetInt(pa_.GetInt() + 1)
                case script.ValueTypeFloat:
                    pa_.SetFloat(pa_.GetFloat() + 1)
                case script.ValueTypeInterface:
                    switch v := pa_.GetInterface().(type) {
                    case script.Int64:
                        pa_.SetInterface(v + 1)
                    case script.Float64:
                        pa_.SetInterface(v + 1)
                    }
                }
            case opcode.Dec:
                switch pa_.GetType() {
                case script.ValueTypeInt:
                    pa_.SetInt(pa_.GetInt() - 1)
                case script.ValueTypeFloat:
                    pa_.SetFloat(pa_.GetFloat() - 1)
                case script.ValueTypeInterface:
                    switch v := pa_.GetInterface().(type) {
                    case script.Int64:
                        pa_.SetInterface(v - 1)
                    case script.Float64:
                        pa_.SetInterface(v - 1)
                    }
                }
            case opcode.Neg:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    pa_.SetInt(-pb_.GetInt())
                case script.ValueTypeFloat:
                    pa_.SetFloat(-pb_.GetFloat())
                case script.ValueTypeInterface:
                    switch v := pb_.GetInterface().(type) {
                    case script.Int64:
                        pa_.SetInterface(-v)
                    case script.Float64:
                        pa_.SetInterface(-v)
                    }
                }
            }
        case opcode.Logic:
            switch il.Code {
            case opcode.Less:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() < pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) < pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) < vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) < vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() < pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() < script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) < vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) < script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ < vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) < vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ < script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) < script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ < vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ < script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ < script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ < script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("<").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.LessOrEqual:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() <= pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) <= pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) <= vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) <= vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() <= pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() <= script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) <= vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) <= script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ <= vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) <= vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ <= script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) <= script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ <= vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ <= script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ <= script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ <= script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet("<=").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Great:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() > pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) > pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) > vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) > vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() > pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() > script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) > vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) > script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ > vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) > vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ > script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) > script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ > vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ > script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ > script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ > script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet(">").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.GreateOrEqual:
                switch pb_.GetType() {
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() >= pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) >= pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) >= vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) >= vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() >= pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() >= script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) >= vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) >= script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ >= vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) >= vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ >= script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) >= script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ >= vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ >= script.Float64(vc_))
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ >= script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ >= script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        pa_.Set(vb_.ScriptGet(">=").GetFunction().Invoke(vb_, pb_.Get()))
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.Equal:
                switch pb_.GetType() {
                case script.ValueTypeBool:
                    switch pc_.GetType() {
                    case script.ValueTypeBool:
                        pa_.SetBool(pb_.GetBool() == pc_.GetBool())
                    default:
                        pa_.SetBool(false)
                    }
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() == pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) == pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) == vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) == vc_)
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() == pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() == script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) == vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) == script.Float64(vc_))
                        default:
                            panic("")
                        }
                    default:
                        panic("")
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ == vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) == vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ == script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) == script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ == vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ == script.Float64(vc_))
                            case script.Object:
                                pa_.SetBool(false)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ == script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ == script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        if vb_.GetScriptTypeId() != script.ScriptTypeNull {
                            switch pc_.GetType() {
                            case script.ValueTypeInterface:
                                switch vc_ := pc_.GetInterface().(type) {
                                case script.String:
                                    if vc_.GetScriptTypeId() == script.ScriptTypeNull {
                                        pa_.SetBool(false)
                                    } else {
                                        pa_.SetBool(vb_ == vc_)
                                    }
                                case script.Object:
                                    if vc_.GetScriptTypeId() == script.ScriptTypeNull {
                                        pa_.SetBool(false)
                                    } else {
                                        pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(vb_, vc_))
                                    }
                                default:
                                    pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(vb_, vc_))
                                }
                            default:
                                pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(vb_, pc_.Get()))
                            }
                        } else {
                            switch pc_.GetType() {
                            case script.ValueTypeInterface:
                                switch vc_ := pc_.GetInterface().(type) {
                                case script.Object:
                                    pa_.SetBool(vc_.GetScriptTypeId() == script.ScriptTypeNull)
                                default:
                                    pa_.SetBool(false)
                                }
                            default:
                                pa_.SetBool(false)
                            }
                        }
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.NotEqual:
                switch pb_.GetType() {
                case script.ValueTypeBool:
                    switch pc_.GetType() {
                    case script.ValueTypeBool:
                        pa_.SetBool(pb_.GetBool() != pc_.GetBool())
                    default:
                        pa_.SetBool(false)
                    }
                case script.ValueTypeInt:
                    switch pc_.GetType() {
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetInt() != pc_.GetInt())
                    case script.ValueTypeFloat:
                        pa_.SetBool(script.Float(pb_.GetInt()) != pc_.GetFloat())
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Int64:
                            pa_.SetBool(script.Int64(pb_.GetInt()) != vc_)
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetInt()) != vc_)
                        default:
                            pa_.SetBool(true)
                        }
                    default:
                        pa_.SetBool(true)
                    }
                case script.ValueTypeFloat:
                    switch pc_.GetType() {
                    case script.ValueTypeFloat:
                        pa_.SetBool(pb_.GetFloat() != pc_.GetFloat())
                    case script.ValueTypeInt:
                        pa_.SetBool(pb_.GetFloat() != script.Float(pc_.GetInt()))
                    case script.ValueTypeInterface:
                        switch vc_ := pc_.GetInterface().(type) {
                        case script.Float64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) != vc_)
                        case script.Int64:
                            pa_.SetBool(script.Float64(pb_.GetFloat()) != script.Float64(vc_))
                        default:
                            pa_.SetBool(true)
                        }
                    default:
                        pa_.SetBool(true)
                    }
                case script.ValueTypeInterface:
                    switch vb_ := pb_.GetInterface().(type) {
                    case script.Int64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Int64:
                                pa_.SetBool(vb_ != vc_)
                            case script.Float64:
                                pa_.SetBool(script.Float64(vb_) != vc_)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ != script.Int64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(script.Float64(vb_) != script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Float64:
                        switch pc_.GetType() {
                        case script.ValueTypeInterface:
                            switch vc_ := pc_.GetInterface().(type) {
                            case script.Float64:
                                pa_.SetBool(vb_ != vc_)
                            case script.Int64:
                                pa_.SetBool(vb_ != script.Float64(vc_))
                            case script.Object:
                                pa_.SetBool(false)
                            default:
                                panic("")
                            }
                        case script.ValueTypeInt:
                            pa_.SetBool(vb_ != script.Float64(pc_.GetInt()))
                        case script.ValueTypeFloat:
                            pa_.SetBool(vb_ != script.Float64(pc_.GetFloat()))
                        default:
                            panic("")
                        }
                    case script.Object:
                        if vb_.GetScriptTypeId() != script.ScriptTypeNull {
                            switch pc_.GetType() {
                            case script.ValueTypeInterface:
                                switch vc_ := pc_.GetInterface().(type) {
                                case script.String:
                                    if vc_.GetScriptTypeId() != script.ScriptTypeNull {
                                        pa_.SetBool(vb_ != vc_)
                                    } else {
                                        pa_.SetBool(true)
                                    }
                                case script.Object:
                                    if vc_.GetScriptTypeId() != script.ScriptTypeNull {
                                        pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(vb_, vc_))
                                    } else {
                                        pa_.SetBool(true)
                                    }
                                default:
                                    pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(vb_, vc_))
                                }
                            default:
                                pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(vb_, pc_.Get()))
                            }
                        } else {
                            switch pc_.GetType() {
                            case script.ValueTypeInterface:
                                switch vc_ := pc_.GetInterface().(type) {
                                case script.Object:
                                    pa_.SetBool(vc_.GetScriptTypeId() != script.ScriptTypeNull)
                                default:
                                    pa_.SetBool(false)
                                }
                            default:
                                pa_.SetBool(false)
                            }
                        }
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.LogicNot:
                pa_.SetBool(!pb_.GetBool())
            case opcode.LogicAnd:
                pa_.SetBool(pb_.GetBool() && pc_.GetBool())
            case opcode.LogicOr:
                pa_.SetBool(pb_.GetBool() || pc_.GetBool())
            }
        case opcode.Flow:
            switch il.Code {
            case opcode.Jump:
                pc = int(pb_.GetInt())
                continue
            case opcode.JumpWhenFalse:
                if !pa_.GetBool() {
                    pc = int(pb_.GetInt())
                    continue
                }
            case opcode.Call:
                switch pa_.GetType() {
                case script.ValueTypeInterface:
                    switch callFunc := pa_.GetInterface().(type) {
                    case script.Function:
                        regStart := pb_.GetInt()
                        count := pc_.GetInt()
                        switch runtimeFunc := callFunc.GetRuntimeFunction().(type) {
                        case runtime.Function:
                            registers = context.PushRegisters(regStart,
                                runtimeFunc.GetMaxRegisterCount()+
                                    len(runtimeFunc.GetArguments())+
                                    len(runtimeFunc.GetLocalVars())+
                                    2)
                            impl.invoke(callFunc)
                            context.PopRegisters()
                        case native.Function:
                            impl.currentPC = pc
                            impl.currentRegisters = registers[:regStart]
                            args := registers[regStart+2 : regStart+2+count]
                            context.PushRegisters(regStart, 1)
                            registers[regStart] = runtimeFunc.NativeCall(registers[regStart+1].Get(), value.ToInterfaceSlice(args)...)
                            context.PopRegisters()
                        }
                    default:
                        panic("")
                    }
                default:
                    panic("")
                }
            case opcode.NewCall:
                cf := pa_.GetInterface().(script.Function)
                regStart := pb_.GetInt()
                count := pc_.GetInt()
                switch rtFunc := cf.GetRuntimeFunction().(type) {
                case runtime.Function:
                    registers = context.PushRegisters(regStart,
                        rtFunc.GetMaxRegisterCount()+
                            len(rtFunc.GetArguments())+
                            len(rtFunc.GetLocalVars())+2)
                    obj := context.NewScriptObject(0)
                    obj.(runtime.Object).SetPrototype(script.InterfaceToValue(cf))
                    registers[regStart+1].SetInterface(obj)
                    impl.invoke(cf)
                    registers[regStart].SetInterface(obj)
                    context.PopRegisters()
                case native.Type:
                    impl.currentPC = pc
                    impl.currentRegisters = registers[:regStart]
                    args := registers[regStart+2 : regStart+2+count]
                    context.PushRegisters(regStart, 1)
                    registers[regStart].SetInterface(rtFunc.New(value.ToInterfaceSlice(args)...))
                    context.PopRegisters()
                case native.Function:
                    panic("can not new with native function")
                }
            case opcode.Ret:
                return
            }
        }

        pc++
    }
}

func NewScriptInterpreter(owner, context interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
        context:       context,
    }
}
