package parser

import (
    "container/list"
)

type (
    MapParser interface {
        ParseMap(interface{}, *list.Element) *list.Element
    }

    ArgListParser interface {
        ParseArgList(interface{}, *list.Element) *list.Element
    }

    ExpressionParser interface {
        ParseExpression(*list.Element) (interface{}, *list.Element)
    }

    FunctionParser interface {
        ParseFunction(interface{}, *list.Element) *list.Element
    }

    ClassParser interface {
        ParseClass(interface{}, *list.Element) *list.Element
    }
)
