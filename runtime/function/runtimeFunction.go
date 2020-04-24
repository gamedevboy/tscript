package function

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"tklibs/script"
	"tklibs/script/debug"
	"tklibs/script/instruction"
	"tklibs/script/opcode"
	"tklibs/script/runtime/runtime_t"
	"tklibs/script/type/function"
)

type Component struct {
	script.ComponentType
	instructions     []instruction.Instruction
	debugInfos       []debug.Info
	arguments        []string
	localVars        []string
	refVars          []string
	members          []string
	name             string
	sourceNames      []string
	isScope          bool
	captureThis      bool
	maxRegisterCount int
	assembly         interface{}
	functions        sync.Map
}

func (impl *Component) RegisterFunction(f uintptr) {
	impl.functions.Store(^uintptr(f), struct{}{})
}

func (impl *Component) UnregisterFunction(f uintptr) {
	impl.functions.Delete(uintptr(f))
}

func (impl *Component) GetAssembly() interface{} {
	return impl.assembly
}

func (impl *Component) CopyFrom(src runtime_t.Function) {
	impl.instructions = src.GetInstructionList()
	impl.debugInfos = src.GetDebugInfoList()
	impl.arguments = src.GetArguments()
	impl.localVars = src.GetLocalVars()
	impl.refVars = src.GetRefVars()
	impl.members = src.GetMembers()
	impl.sourceNames = src.GetSourceNames()
	impl.isScope = src.IsScope()
	impl.captureThis = src.IsCaptureThis()
	impl.maxRegisterCount = src.GetMaxRegisterCount()

	impl.functions.Range(func(key, val interface{}) bool {
		ptr := (*function.Component)(unsafe.Pointer(^(key.(uintptr))))
		ptr.Reload()
		return true
	})
}

func (impl *Component) GetMaxRegisterCount() int {
	return impl.maxRegisterCount
}

func (impl *Component) GetDebugInfoList() []debug.Info {
	return impl.debugInfos
}

var _ runtime_t.Function = &Component{}

func (impl *Component) IsScope() bool {
	return impl.isScope
}

func (impl *Component) IsCaptureThis() bool {
	return impl.captureThis
}

func (impl *Component) GetInstructionList() []instruction.Instruction {
	return impl.instructions
}

func (impl *Component) GetArguments() []string {
	return impl.arguments
}

func (impl *Component) GetLocalVars() []string {
	return impl.localVars
}

func (impl *Component) GetRefVars() []string {
	return impl.refVars
}

func (impl *Component) GetMembers() []string {
	return impl.members
}

func (impl *Component) GetName() string {
	return impl.name
}

func (impl *Component) GetSourceNames() []string {
	return impl.sourceNames
}

func (impl *Component) String() string {
	return fmt.Sprintf("Func<%s>", impl.name)
}

func (impl *Component) DumpString() string {
	builder := strings.Builder{}
	debugPcIndex := 0
	debugInfoList := impl.GetDebugInfoList()
	for i, il := range impl.instructions {
		instStr := ""

		tc := il.Type & 3
		tb := (il.Type >> 2) & 3

		ra, rb, rc := "", "", ""
		if il.A < 0 {
			ra = fmt.Sprintf("$%v", -il.A-1)
		} else {
			ra = fmt.Sprintf("[%v]", il.A)
		}

		switch tc {
		case opcode.Register:
			rc = fmt.Sprintf("[%v]", il.C)
		case opcode.Reference:
			rc = fmt.Sprintf("$%v", il.C)
		case opcode.Integer:
			rc = fmt.Sprintf("%v", il.C)
		}

		if tc == opcode.None {
			switch tb {
			case opcode.Register:
				rb = fmt.Sprintf("[%v]", il.B)
			case opcode.Reference:
				rb = fmt.Sprintf("$%v", il.B)
			case opcode.Integer:
				rb = fmt.Sprintf("%v", il.GetABx().B)
			case opcode.None:
				rb = fmt.Sprintf("%v", il.GetABm().B)
			}
		} else {
			switch tb {
			case opcode.Register:
				rb = fmt.Sprintf("[%v]", il.B)
			case opcode.Reference:
				rb = fmt.Sprintf("$%v", il.B)
			case opcode.Integer:
				rb = fmt.Sprintf("%v", il.B)
			}
		}

		switch il.Type >> 4 {
		case opcode.Nop:
			{
				instStr = "NOP"
			}
		case opcode.Memory:
			switch il.Code {
			case opcode.Move:
				instStr = fmt.Sprintf("MOVE \t%v, \t%v", ra, rb)
			case opcode.Object:
				instStr = fmt.Sprintf("OBJ \t%v, \t%v", ra, rb)
			case opcode.Array:
				instStr = fmt.Sprintf("ARRAY \t%v, \t%v", ra, rb)
			case opcode.LoadField:
				instStr = fmt.Sprintf("LDFLD \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.StoreField:
				instStr = fmt.Sprintf("STFLD \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.LoadElement:
				instStr = fmt.Sprintf("LDELE \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.StoreElement:
				instStr = fmt.Sprintf("STELE \t%v, \t%v, \t%v", ra, rb, rc)
			}
		case opcode.Const:
			switch il.Code {
			case opcode.Load:
				instStr = fmt.Sprintf("LOAD \t%v, \t%v", ra, rb)
			case opcode.LoadNil:
				instStr = fmt.Sprintf("LDNIL \t%v", ra)
			case opcode.LoadFunc:
				instStr = fmt.Sprintf("LDFN \t%v, \t%v", ra, rb)
			case opcode.LoadBool:
				instStr = fmt.Sprintf("LDBOOL \t%v, \t%v", ra, rb)
			}
		case opcode.Math:
			switch il.Code {
			case opcode.Add:
				instStr = fmt.Sprintf("ADD \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Sub:
				instStr = fmt.Sprintf("SUB \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Mul:
				instStr = fmt.Sprintf("MUL \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Div:
				instStr = fmt.Sprintf("DIV \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Inc:
				instStr = fmt.Sprintf("INC \t%v", ra)
			case opcode.Dec:
				instStr = fmt.Sprintf("DEC \t%v", ra)
			case opcode.Neg:
				instStr = fmt.Sprintf("NEG \t%v, \t%v", ra, rb)
			case opcode.Rem:
				instStr = fmt.Sprintf("REM \t%v, \t%v, \t%v", ra, rb, rc)
			}
		case opcode.Logic:
			switch il.Code {
			case opcode.Equal:
				instStr = fmt.Sprintf("EQ \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.NotEqual:
				instStr = fmt.Sprintf("NEQ \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Less:
				instStr = fmt.Sprintf("LESS \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Great:
				instStr = fmt.Sprintf("GREAT \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.LessOrEqual:
				instStr = fmt.Sprintf("LEQ \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.GreateOrEqual:
				instStr = fmt.Sprintf("GEQ \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.LogicOr:
				instStr = fmt.Sprintf("LOR \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.LogicAnd:
				instStr = fmt.Sprintf("LAND \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.LogicNot:
				instStr = fmt.Sprintf("LNOT \t%v, \t%v", ra, rb)
			}
		case opcode.Flow:
			switch il.Code {
			case opcode.Call:
				instStr = fmt.Sprintf("CALL \t%v, \t%v, \t%v", ra, rb, rc)
			case opcode.Jump:
				if il.GetABx().B > 0 {
					instStr = fmt.Sprintf("JEZ \t%v, \t%v", ra, il.GetABx().B)
				} else {
					instStr = fmt.Sprintf("JNZ \t%v, \t%v", ra, -il.GetABx().B)
				}
			case opcode.JumpTo:
				instStr = fmt.Sprintf("JMP \t%v", rb)
			case opcode.Ret:
				instStr = "RET"
			}
		}

		if len(debugInfoList) > 0 {
			for debugPcIndex < len(debugInfoList) && uint32(i) > debugInfoList[debugPcIndex].PC {
				debugPcIndex++
			}

			if debugPcIndex > 0 {
				debugPcIndex--
			}
		}

		builder.WriteString(fmt.Sprintf("%v:%v \t\t %v\n", i, debugInfoList[debugPcIndex].Line, instStr))
	}

	return builder.String()

}

func NewFunctionComponent(owner, assembly interface{}, instructionCount int, debugInfoCount int, name string, sourceFiles []string,
	localVars,
	arguments, refVars,
	members []string, isScope bool, captureThis bool, maxRegisterCount uint32) *Component {
	return &Component{
		ComponentType:    script.MakeComponentType(owner),
		assembly:         assembly,
		instructions:     make([]instruction.Instruction, instructionCount),
		debugInfos:       make([]debug.Info, debugInfoCount),
		arguments:        arguments,
		localVars:        localVars,
		refVars:          refVars,
		members:          members,
		name:             name,
		sourceNames:      sourceFiles,
		isScope:          isScope,
		captureThis:      captureThis,
		maxRegisterCount: int(maxRegisterCount),
	}
}
