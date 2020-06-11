package ast

import "container/list"

type Statement interface {
    Node
    Compile(interface{}) *list.Element
}
