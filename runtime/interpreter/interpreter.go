package interpreter

import (
	"fmt"
	"strings"
	"unsafe"

	"tklibs/script"
	"tklibs/script/instruction"
	"tklibs/script/opcode"
	"tklibs/script/runtime"
	"tklibs/script/runtime/runtime_t"
	"tklibs/script/runtime/scope"
	"tklibs/script/runtime/stack/frame"
	"tklibs/script/runtime/util"
	"tklibs/script/type/function"
	"tklibs/script/value"
)

type Component struct {
	script.ComponentType
	context          runtime.ScriptContext
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

func (impl *Component) InvokeNew(function, context interface{}, args ...interface{}) interface{} {
	this := impl.context.NewScriptObject(0)
	this.(runtime.Object).SetPrototype(script.InterfaceToValue(function))
	sf := function.(script.Function)
	switch _func := sf.GetRuntimeFunction().(type) {
	case runtime_t.Function:
		if len(args) < len(_func.GetArguments()) {
			panic("") // todo not enough arguments
		}
		context := context.(runtime.ScriptContext)
		context.PushRegisters(0, _func.GetMaxRegisterCount()+len(_func.GetLocalVars())+len(_func.GetArguments())+2)
		defer context.PopRegisters()
		registers := context.GetRegisters()
		defer registers[0].SetNull()
		registers[1].Set(this)
		for i := range _func.GetArguments() {
			registers[2+i].Set(args[i])
		}
		impl.invoke(function.(script.Function))
		ret := registers[0].Get()
		return ret
	case runtime_t.NativeFunction:
		ret, _ := _func.NativeCall(context, this, args...)
		return ret
	default:
		panic("")
	}
}

func (impl *Component) InvokeFunction(function, context, this interface{}, args ...interface{}) interface{} {
	sf := function.(script.Function)
	switch _func := sf.GetRuntimeFunction().(type) {
	case runtime_t.Function:
		if len(args) < len(_func.GetArguments()) {
			panic(fmt.Sprintf("not enough arguments,get:%d excepted:%d", len(args), len(_func.GetArguments()))) // todo not enough arguments
		}
		context := context.(runtime.ScriptContext)
		context.PushRegisters(0, _func.GetMaxRegisterCount()+len(_func.GetLocalVars())+len(_func.GetArguments())+2)
		defer context.PopRegisters()
		registers := context.GetRegisters()
		defer registers[0].SetNull()
		registers[1].Set(this)
		for i := range _func.GetArguments() {
			registers[2+i].Set(args[i])
		}
		impl.invoke(function.(script.Function))
		ret := registers[0].Get()
		return ret
	case runtime_t.NativeFunction:
		ret, _ := _func.NativeCall(context, this, args...)
		return ret
	default:
		panic("")
	}
}

//noinspection GoNilness
func (impl *Component) invoke(sf script.Function) (exception interface{}) {
	_func := sf.GetScriptRuntimeFunction()

	exception = nil

	context := impl.context
	registers := context.GetRegisters()
	registers[0].SetNull()

	defer func() {
		regPtr := uintptr(unsafe.Pointer(&registers[1]))

		len := _func.GetMaxRegisterCount() +
			len(_func.GetArguments()) +
			len(_func.GetLocalVars()) +
			2

		for i := 1; i < len; i++ {
			*(**interface{})(unsafe.Pointer(regPtr)) = nil
			regPtr += uintptr(8)
		}
	}()

	if _func.IsCaptureThis() {
		registers[1] = sf.GetThis()
	}

	pc := 0

	defer func() {
		if err := recover(); err != nil {
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

			switch e := err.(type) {
			case script.Error:
				panic(script.MakeError(fileName, line, "%v @ %v: %v", e, fileName, line))
			case script.ScriptException:
				exception = e.GetException()
			default:
				panic(script.MakeError(fileName, line, "script runtime error: [%v] @ %v:%v in %v", err, fileName, line, _func.GetName()))
			}
		}
	}()

	instList := _func.GetInstructionList()
	instCount := len(instList)
	if instCount == 0 {
		registers[0].SetNull()
		return
	}

	if _func.IsScope() {
		context.PushScope(scope.NewScope(nil, sf, registers[2:], registers[2+len(_func.GetArguments()):]))
		defer freeScope(context)
	}

	var vb, vc script.Value
	var pa_, pb_, pc_ *script.Value

	ilStart := uintptr(unsafe.Pointer(&instList[0]))
	ilPtr := ilStart
	il := (*instruction.Instruction)(unsafe.Pointer(ilPtr))

vm_loop:
	for pc < instCount {
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
			vb.SetInt(script.Int(il.GetABx().B))
			pb_ = &vb
		case opcode.None:
			b := il.GetABm().B
			vb.SetFloat(script.Float(b))
			pb_ = &vb
		}

		switch _type {
		case opcode.Memory:
			switch il.Code {
			case opcode.Move:
				*pa_ = *pb_
			case opcode.LoadField:
				index := pc_.GetInt()
				obj := pb_.Get()
				if index > -1 {
					*pa_ = sf.GetFieldByMemberIndex(obj, index)
				} else {
					if obj != nil && obj != script.Null {
						*pa_ = sf.GetFieldByMemberIndex(obj, -index-1)
					} else {
						*pa_ = script.NullValue
					}
				}
			case opcode.StoreField:
				index := pb_.GetInt()
				target := pa_.Get()
				if index > -1 {
					sf.SetFieldByMemberIndex(target, index, *pc_)
				} else {
					if target != nil && target != script.Null {
						sf.SetFieldByMemberIndex(target, -index-1, *pc_)
					}
				}
			case opcode.LoadElement:
				switch pb_.GetType() {
				case script.ValueTypeInterface:
					switch target := pb_.GetInterface().(type) {
					case script.Map:
						pa_.Set(target.Get(pc_.Get()))
					case script.Array:
						*pa_ = target.GetElement(pc_.ToInt())
					case script.Object:
						*pa_ = target.ScriptGet(string(util.ToScriptString(context, pc_.Get())))
					}
				default:
					panic("")
				}
			case opcode.StoreElement:
				switch pa_.GetType() {
				case script.ValueTypeInterface:
					switch target := pa_.GetInterface().(type) {
					case script.Map:
						target.Set(pb_.Get(), pc_.Get())
					case script.Array:
						target.SetElement(pb_.ToInt(), *pc_)
					case script.Object:
						target.ScriptSet(string(util.ToScriptString(context, pb_.Get())), *pc_)
					}
				default:
					panic("")
				}
			case opcode.Object:
				pa_.SetInterface(context.NewScriptObject(int(pb_.GetInt())))
			case opcode.Array:
				pa_.SetInterface(context.NewScriptArray(int(pb_.GetInt())))
			}
		case opcode.Const:
			switch il.Code {
			case opcode.Load:
				assembly := _func.GetAssembly()
				index := pb_.GetInt()
				_t := int(index & 3)
				index = index >> 2
				switch _t {
				case opcode.ConstInt64:
					pa_.SetInt64(assembly.(script.Assembly).GetIntConstPool().Get(int(index)).(script.Int64))
				case opcode.ConstFloat64:
					pa_.SetFloat64(assembly.(script.Assembly).GetFloatConstPool().Get(int(index)).(script.Float64))
				case opcode.ConstString:
					pa_.SetInterface(script.String(assembly.(script.Assembly).GetStringConstPool().Get(int(index)).(string)))
				case opcode.ConstBool:
					pa_.SetBool(index != 0)
				}
			case opcode.LoadFunc:
				metaIndex := pb_.GetInt()
				f := &struct {
					*function.Component
				}{}
				f.Component = function.NewScriptFunction(f, _func.GetAssembly().(script.Assembly).
					GetFunctionByMetaIndex(metaIndex),
					context)
				rf := f.GetRuntimeFunction().(runtime_t.Function)
				f.Init()
				if rf.IsCaptureThis() {
					f.SetThis(registers[1])
				}
				pa_.SetInterface(f)
			case opcode.LoadNil:
				*pa_ = script.NullValue
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
					case script.String:
						pa_.Set(context.GetStringPrototype().(script.Object).ScriptGet("+").GetFunction().Invoke(context, vb_, util.ToScriptString(context, pc_.Get())))
					case script.Object:
						fn := vb_.ScriptGet("+")
						if fn.IsNull() || fn.GetPointerType() != script.InterfaceTypeFunction {
							panic("Can't find '+' operator")
						}
						pa_.Set(fn.GetFunction().Invoke(context, vb_, pc_.Get()))
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
						pa_.Set(vb_.ScriptGet("-").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet("*").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet("/").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet("%").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet("<").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet("<=").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet(">").GetFunction().Invoke(context, vb_, pb_.Get()))
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
						pa_.Set(vb_.ScriptGet(">=").GetFunction().Invoke(context, vb_, pb_.Get()))
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
							if pc_.IsNull() || vc_ == script.Null {
								pa_.SetBool(false)
							} else {
								panic("")
							}
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
					case script.String:
						if pc_.IsNull() {
							pa_.SetBool(false)
						} else {
							switch vc_ := pc_.GetInterface().(type) {
							case script.String:
								pa_.SetBool(strings.Compare(string(vb_), string(vc_)) == 0)
							default:
								if vc_ == script.Null {
									pa_.SetBool(false)
								} else {
									panic("")
								}
							}
						}
					case script.Object:
						if vb_.GetScriptTypeId() != script.ScriptTypeNull {
							switch pc_.GetType() {
							case script.ValueTypeInterface:
								switch vc_ := pc_.GetInterface().(type) {
								case script.Object:
									if vc_.GetScriptTypeId() == script.ScriptTypeNull {
										pa_.SetBool(false)
									} else {
										pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(context, vb_, vc_))
									}
								default:
									pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(context, vb_, vc_))
								}
							default:
								pa_.Set(vb_.ScriptGet("==").GetFunction().Invoke(context, vb_, pc_.Get()))
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
					case script.String:
						if pc_.IsNull() {
							pa_.SetBool(true)
						} else {
							switch vc_ := pc_.GetInterface().(type) {
							case script.String:
								pa_.SetBool(strings.Compare(string(vb_), string(vc_)) != 0)
							default:
								if vc_ == script.Null {
									pa_.SetBool(true)
								} else {
									panic("")
								}
							}
						}
					case script.Object:
						if vb_.GetScriptTypeId() != script.ScriptTypeNull {
							switch pc_.GetType() {
							case script.ValueTypeInterface:
								switch vc_ := pc_.GetInterface().(type) {
								case script.Object:
									if vc_.GetScriptTypeId() != script.ScriptTypeNull {
										pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(context, vb_, vc_))
									} else {
										pa_.SetBool(true)
									}
								default:
									pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(context, vb_, vc_))
								}
							default:
								pa_.Set(vb_.ScriptGet("!=").GetFunction().Invoke(context, vb_, pc_.Get()))
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
						panic(fmt.Errorf("Unknow object %v ", vb_))
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
			case opcode.JumpTo:
				pc = int(pb_.GetInt())
				ilPtr = ilStart + uintptr(int(pc)*8)
				il = (*instruction.Instruction)(unsafe.Pointer(ilPtr))
				continue
			case opcode.Jump:
				p := int(pb_.GetInt())
				if p > 0 {
					if pa_.IsNull() || bool(!pa_.GetBool()) {
						pc = p
						ilPtr = ilStart + uintptr(pc*8)
						il = (*instruction.Instruction)(unsafe.Pointer(ilPtr))
						continue
					}
				} else {
					if !pa_.IsNull() && bool(pa_.GetBool()) {
						pc = -p
						ilPtr = ilStart + uintptr(pc*8)
						il = (*instruction.Instruction)(unsafe.Pointer(ilPtr))
						continue
					}
				}
			case opcode.JumpNull:
				p := int(pb_.GetInt())
				if !pa_.IsNull() {
					pc = p
					ilPtr = ilStart + uintptr(pc*8)
					il = (*instruction.Instruction)(unsafe.Pointer(ilPtr))
					continue
				}
			case opcode.Call:
				regStart := pb_.GetInt()
				count := pc_.GetInt()

				isNewCall := regStart&1 != 0
				regStart = regStart >> 1

				if count < 0 {
					count = -count
					start := int(regStart + 1 + count)
					if array, ok := registers[start].Get().(script.Array); ok {
						for i := 0; i < int(array.Len()); i++ {
							registers[start+i] = array.GetElement(script.Int(i))
						}
					}
				}

				switch pa_.GetType() {
				case script.ValueTypeInterface:
					pai := pa_.Interface()

					switch pai.GetType() {
					case script.InterfaceTypeFunction:
						callFunc := pai.GetFunction()

						_frame := frame.NewStackFrame(nil, _func)
						context.PushFrame(_frame)

						if callFunc.IsScriptFunction() {
							runtimeFunc := callFunc.GetScriptRuntimeFunction()
							registers = context.PushRegisters(regStart,
								runtimeFunc.GetMaxRegisterCount()+
									len(runtimeFunc.GetArguments())+
									len(runtimeFunc.GetLocalVars())+
									2)

							var newObj interface{}

							if isNewCall {
								newObj = context.NewScriptObject(0)
								newObj.(runtime.Object).SetPrototype(script.InterfaceToValue(callFunc))
								registers[regStart+1].SetInterface(newObj)
							}

							impl.invoke(callFunc)

							if isNewCall {
								registers[regStart].SetInterface(newObj)
							}

							context.PopRegisters()
						} else {
							impl.currentPC = pc
							impl.currentRegisters = registers[:regStart]
							args := registers[regStart+2 : regStart+2+count]

							context.PushRegisters(regStart, 1)

							switch f := callFunc.(type) {
							case runtime_t.NativeFunction:
								ret, _ := f.NativeCall(context, registers[regStart+1].Get(), value.ToInterfaceSlice(args)...)
								registers[regStart].Set(ret)
							default:
								fc := callFunc.GetNativeRuntimeFunction()
								if fc != nil {
									ret, _ := fc.NativeCall(context, registers[regStart+1].Get(), value.ToInterfaceSlice(args)...)
									registers[regStart].Set(ret)
								} else {
									panic("Invalid native function call")
								}
							}
							context.PopRegisters()
						}
						frame.FreeStackFrame(context.PopFrame().(*frame.Component))
					case script.InterfaceTypeAny:
						if pai.IsNull() {
							panic(fmt.Errorf("Cannot call on 'null' value "))
						}

						cm := pai.Get().(script.Object).ScriptGet("()")

						if cm.IsNull() {
							panic(fmt.Errorf("Cannot find call operator in object "))
						}

						callMethod := cm.Get()
						if callFunc, ok := callMethod.(script.Function); ok {
							switch runtimeFunc := callFunc.GetRuntimeFunction().(type) {
							case runtime_t.Function:
								registers = context.PushRegisters(regStart,
									runtimeFunc.GetMaxRegisterCount()+
										len(runtimeFunc.GetArguments())+
										len(runtimeFunc.GetLocalVars())+
										2)
								impl.invoke(callFunc)
								context.PopRegisters()
							case runtime_t.NativeFunction:
								impl.currentRegisters = registers[:regStart]
								args := registers[regStart+2 : regStart+2+count]
								context.PushRegisters(regStart, 1)
								ret, _ := runtimeFunc.NativeCall(context, registers[regStart+1].Get(), value.ToInterfaceSlice(args)...)
								registers[regStart].Set(ret)
								context.PopRegisters()
							}
						}
					}
				default:
					panic("")
				}
			case opcode.Ret:
				break vm_loop
			}
		case opcode.Bit:
			switch il.Code {
			case opcode.Or:
				switch pb_.GetType() {
				case script.ValueTypeInt:
					switch pc_.GetType() {
					case script.ValueTypeInt:
						pa_.SetInt(pb_.GetInt() | pc_.GetInt())
					case script.ValueTypeInterface:
						switch vc_ := pc_.GetInterface().(type) {
						case script.Int64:
							pa_.SetInt64(script.Int64(pb_.GetInt()) | vc_)
						default:
							panic("not support |")
						}
					default:
						panic("not support |")
					}
				case script.ValueTypeBool:
					switch pc_.GetType() {
					case script.ValueTypeBool:
						pa_.SetBool(pb_.GetBool() || pc_.GetBool())
					default:
						panic("not support |")
					}
				case script.ValueTypeInterface:
					switch vb_ := pb_.GetInterface().(type) {
					case script.Int64:
						switch pc_.GetType() {
						case script.ValueTypeInterface:
							switch vc_ := pc_.GetInterface().(type) {
							case script.Int64:
								pa_.SetInt64(vb_ | vc_)
							default:
								panic("not support |")
							}
						case script.ValueTypeInt:
							pa_.SetInt64(vb_ | script.Int64(pc_.GetInt()))
						default:
							panic("not support |")
						}
					default:
						panic("not support |")
					}
				default:
					panic(fmt.Errorf("Add can not support: %v: %v ", pb_.GetType(), pb_.Get()))
				}
			case opcode.And:
				switch pb_.GetType() {
				case script.ValueTypeInt:
					switch pc_.GetType() {
					case script.ValueTypeInt:
						pa_.SetInt(pb_.GetInt() & pc_.GetInt())
					case script.ValueTypeInterface:
						switch vc_ := pc_.GetInterface().(type) {
						case script.Int64:
							pa_.SetInt64(script.Int64(pb_.GetInt()) & vc_)
						default:
							panic("not support |")
						}
					default:
						panic("not support |")
					}
				case script.ValueTypeBool:
					switch pc_.GetType() {
					case script.ValueTypeBool:
						pa_.SetBool(pb_.GetBool() && pc_.GetBool())
					default:
						panic("not support |")
					}
				case script.ValueTypeInterface:
					switch vb_ := pb_.GetInterface().(type) {
					case script.Int64:
						switch pc_.GetType() {
						case script.ValueTypeInterface:
							switch vc_ := pc_.GetInterface().(type) {
							case script.Int64:
								pa_.SetInt64(vb_ & vc_)
							default:
								panic("not support |")
							}
						case script.ValueTypeInt:
							pa_.SetInt64(vb_ & script.Int64(pc_.GetInt()))
						default:
							panic("not support |")
						}
					default:
						panic("not support |")
					}
				default:
					panic(fmt.Errorf("Add can not support: %v: %v ", pb_.GetType(), pb_.Get()))
				}
			}
		}

		ilPtr += uintptr(8)
		il = (*instruction.Instruction)(unsafe.Pointer(ilPtr))
		pc++
	}

	return
}

func freeScope(context runtime.ScriptContext) {
	scope.FreeScope(context.PopScope().(*scope.Component))
}

func NewScriptInterpreter(owner, context interface{}) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
		context:       context.(runtime.ScriptContext),
	}
}
