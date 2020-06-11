package parser

import (
	"tklibs/script/compiler/token"
)

type (
	BlockParser interface {
		ParseBlock(interface{}, *token.Iterator) *token.Iterator
		ParseStatement(element *token.Iterator) (interface{}, *token.Iterator)
	}

	DeclParser interface {
		ParseDecl(interface{}, *token.Iterator) *token.Iterator
	}

	IfParser interface {
		ParseIf(interface{}, *token.Iterator) *token.Iterator
	}

	WhileParser interface {
		ParseWhile(interface{}, *token.Iterator) *token.Iterator
	}

	ForParser interface {
		ParseFor(interface{}, *token.Iterator) *token.Iterator
	}

	ForeachParser interface {
		ParseForeach(interface{}, *token.Iterator) *token.Iterator
	}

	SwitchParser interface {
		ParseSwitch(interface{}, *token.Iterator) *token.Iterator
	}
)
