package arglist

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler/ast/expression"
)

type Component struct {
    script.ComponentType
    expressionList list.List
}

func (impl *Component) GetExpressionList() *list.List {
    return &impl.expressionList
}

func (impl *Component) String() string {
    r := ""
    for it := impl.expressionList.Front(); it != nil; it = it.Next() {
        r += fmt.Sprint(it.Value)
        if it.Next() != nil {
            r += fmt.Sprint(",")
        }
    }

    return r
}

var _ expression.ArgList = &Component{}

func NewArgList(owner interface{}) *Component {
    return &Component{ComponentType: script.MakeComponentType(owner)}
}
