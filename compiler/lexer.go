package compiler

import (
	"container/list"
)

type Lexer interface {
	ParseFromRunes(content []rune) *list.List
}

