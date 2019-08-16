package arglist

import (
    "container/list"
    "fmt"

    "tklibs/script"
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

func NewArgList(owner interface{}) *Component {
    return &Component{ComponentType: script.MakeComponentType(owner)}
}
