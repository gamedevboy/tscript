package trivia

import (
	"container/list"
	"strings"

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
	content string
}

func (c *Component) SetContent(content string) {
	c.content = content
}

func (c *Component) GetContent() string {
	return c.content
}

func (c *Component) Format(ident int, formatBuilder *strings.Builder) {
	formatBuilder.WriteString(c.content)
}

func (c *Component) Compile(f interface{}) *list.Element {
	_func := f.(compiler.Function)
	return _func.GetInstructionList().Back()
}
