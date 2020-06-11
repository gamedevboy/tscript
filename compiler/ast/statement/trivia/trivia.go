package trivia

import (
	"container/list"

	"tklibs/script"
	"tklibs/script/compiler"
	"tklibs/script/compiler/ast"
	"tklibs/script/compiler/debug"
)

var _ ast.Statement = &Component{}

func NewTrivia(owner interface{}) *Component {
	return &Component{
		ComponentType: script.MakeComponentType(owner),
	}
}

type Component struct {
	debug.Component
	script.ComponentType
}

func (c *Component) String() string {
	panic("implement me")
}

func (c *Component) Compile(f interface{}) *list.Element {
	_func := f.(compiler.Function)
	return _func.GetInstructionList().Back()
}
