package _case

import (
    "container/list"
    "strings"

    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/debug"
)

type Component struct {
    debug.Component
    script.ComponentType

    value interface{}
    block interface{}
}

func (impl *Component) Format(ident int, formatBuilder *strings.Builder) {
    panic("implement me")
}

func (impl *Component) Compile(i interface{}) *list.Element {
    panic("implement me")
}

var _ statement.Case = &Component{}

func NewCase(owner interface{}) *Component {
    return &Component{
        ComponentType: script.MakeComponentType(owner),
    }
}

func (impl *Component) GetValue() interface{} {
    return impl.value
}

func (impl *Component) SetValue(value interface{}) {
    impl.value = value
}

func (impl *Component) GetBlock() interface{} {
    return impl.block
}

func (impl *Component) SetBlock(value interface{}) {
    impl.block = value
}
