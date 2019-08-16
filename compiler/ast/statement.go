package ast

import "container/list"

type Statement interface {
    Compile(interface{}) *list.Element
}
