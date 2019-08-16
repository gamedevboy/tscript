package ast

import "tklibs/script/compiler"

type Expression interface {
    Compile(interface{}, *compiler.Operand) *compiler.Operand
}
