package ast

import "tklibs/script/compiler"

type Expression interface {
    Node
    Compile(interface{}, *compiler.Operand) *compiler.Operand
}
