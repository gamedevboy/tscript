package _const

import (
    "fmt"
    "math"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/opcode"
)

type Component struct {
    script.ComponentType
    value interface{}
}

func (impl *Component) String() string {
    return impl.value.(string)
}

func (impl *Component) GetValue() interface{} {
    return impl.value
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
    _func := f.(compiler.Function)
    asm := _func.GetAssembly().(script.Assembly)

    switch value := impl.value.(type) {
    case script.Int:
        if value < math.MinInt16 || value > math.MaxInt16 {
            if r == nil {
                r = compiler.NewRegisterOperand(_func.AllocRegister(""))
            }
            _func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewIntOperand(value))
            return r
        } else {
            if r == nil {
                return compiler.NewSmallIntOperand(int16(value))
            }
            _func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewIntOperand(value))
        }
    case script.Float:
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(""))
        }
        _func.AddInstructionABm(opcode.Move, opcode.Memory, r, compiler.NewFloatOperand(value))
        return r
    case script.Int64:
        rt := fmt.Sprintf("Int64:%v", value)
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(rt))
        } else {
            if r.Register != nil {
                r.Tag = rt
            }
        }
        _func.AddInstructionABx(opcode.Load, opcode.Const, r, compiler.NewIntOperand(script.Int(asm.GetIntConstPool().Insert(
            value)<<2)+script.Int(opcode.ConstInt64)))
    case script.Float64:
        rt := fmt.Sprintf("Float64:%v", value)
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(rt))
        } else {
            if r.Register != nil {
                r.Tag = rt
            }
        }
        _func.AddInstructionABx(opcode.Load, opcode.Const, r, compiler.NewIntOperand(script.Int(asm.GetFloatConstPool().Insert(
            value)<<2)+script.Int(opcode.ConstFloat64)))
    case script.String:
        switch value {
        case "this":
            if r == nil {
                r = compiler.NewRegisterOperand(&compiler.Register{Index: 1})
            } else {
                _func.AddInstructionABx(opcode.Move, opcode.Memory, r, compiler.NewRegisterOperand(&compiler.Register{Index: 1}))
            }

            return r
        case "null":
            if r == nil {
                r = compiler.NewRegisterOperand(_func.AllocRegister(""))
            }
            _func.AddInstructionABx(opcode.LoadNil, opcode.Const, r, nil)
            return r
        default:
            rt := fmt.Sprintf("String:%v", value)
            if r == nil {
                r = compiler.NewRegisterOperand(_func.AllocRegister(rt))
            } else {
                if r.Register != nil {
                    r.Tag = rt
                }
            }
            _func.AddInstructionABx(opcode.Load, opcode.Const, r, compiler.NewIntOperand(script.Int(asm.GetStringConstPool().Insert(
                string(value))<<2)+script.Int(opcode.ConstString)))
        }
    case script.Bool:
        rt := fmt.Sprintf("Bool:%v", value)
        if r == nil {
            r = compiler.NewRegisterOperand(_func.AllocRegister(rt))
        } else {
            if r.Register != nil {
                r.Register.Tag = rt
            }
        }
        if value {
            _func.AddInstructionABx(opcode.LoadBool, opcode.Const, r, compiler.NewIntOperand(1))
        } else {
            _func.AddInstructionABx(opcode.LoadBool, opcode.Const, r, compiler.NewIntOperand(0))
        }
    }

    return r
}

func NewConst(owner, value interface{}) *Component {
    return &Component{
        script.MakeComponentType(owner),
        value,
    }
}
