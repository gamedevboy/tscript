package parser

import (
    "container/list"
)

type (
    BlockParser interface {
        ParseBlock(interface{}, *list.Element) *list.Element
        ParseStatement(element *list.Element) (interface{}, *list.Element)
    }
    DeclParser interface {
        ParseDecl(interface{}, *list.Element) *list.Element
    }
    IfParser interface {
        ParseIf(interface{}, *list.Element) *list.Element
    }
    WhileParser interface {
        ParseWhile(interface{}, *list.Element) *list.Element
    }
    ForParser interface {
        ParseFor(interface{}, *list.Element) *list.Element
    }
    ForeachParser interface {
        ParseForeach(interface{}, *list.Element) *list.Element
    }
    SwitchParser interface {
        ParseSwitch(interface{}, *list.Element) *list.Element
    }
)
