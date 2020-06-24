package function

import (
	"fmt"
	"strings"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/compiler/debug"
	"tklibs/script/opcode"
)

type Component struct {
	debug.Component
	script.ComponentType
	argList     interface{}
	block       interface{}
	name        string
	metaIndex   uint32
	captureThis bool
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
	if impl.captureThis {
		formatBuilder.WriteString("(")
		impl.argList.(ast.Node).Format(ident, formatBuilder)
		formatBuilder.WriteString(") => {")
		impl.block.(ast.Node).Format(ident+4, formatBuilder)
		formatBuilder.WriteString("}")
	} else {
		formatBuilder.WriteString(fmt.Sprintf("function %v(", impl.name))
		impl.argList.(ast.Node).Format(ident, formatBuilder)
		formatBuilder.WriteString(") {")
		impl.block.(ast.Node).Format(ident+4, formatBuilder)
		formatBuilder.WriteString("}")
	}
}

func (impl *Component) SetCaptureThis(val bool) {
	impl.captureThis = val
}

func (impl *Component) GetCaptureThis() bool {
	return impl.captureThis
}

var _ expression.Function = &Component{}

func (impl *Component) SetMetaIndex(value uint32) {
	impl.metaIndex = value
}

func (impl *Component) GetMetaIndex() uint32 {
	return impl.metaIndex
}

func (impl *Component) GetName() string {
	return impl.name
}

func (impl *Component) SetName(name string) {
	impl.name = name
}

func (impl *Component) GetArgList() interface{} {
	return impl.argList
}

func (impl *Component) SetArgList(arglist interface{}) {
	impl.argList = arglist
}

func (impl *Component) GetBlock() interface{} {
	return impl.block
}

func (impl *Component) SetBlock(bs interface{}) {
	impl.block = bs
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
	_func := f.(compiler.Function)
	if r == nil {
		r = compiler.NewRegisterOperand(_func.AllocRegister(""))
	}

	_func.AddInstructionABx(opcode.LoadFunc, opcode.Const, r, compiler.NewIntOperand(script.Int(impl.metaIndex)))

	return r
}

func NewFunction(owner interface{}) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
	}
}
