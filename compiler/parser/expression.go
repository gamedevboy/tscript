package parser

import (
    "tklibs/script/compiler/token"
)

type (
    ObjectParser interface {
        ParseObject(interface{}, *token.Iterator) *token.Iterator
    }

    ArgListParser interface {
        ParseArgList(interface{}, *token.Iterator) *token.Iterator
    }

    ExpressionParser interface {
        ParseExpression(*token.Iterator) (interface{}, *token.Iterator)
    }

    FunctionParser interface {
        ParseFunction(interface{}, *token.Iterator) *token.Iterator
    }

    ClassParser interface {
        ParseClass(interface{}, *token.Iterator) *token.Iterator
    }
)
