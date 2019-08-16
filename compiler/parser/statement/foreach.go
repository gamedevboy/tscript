package statement

import (
    "container/list"

    "tklibs/script"
)

type ForeachStatementParserComponent struct {
    *script.ComponentType
}

func (parser *ForeachStatementParserComponent) ParseForeach(f interface{}, tokenIt *list.Element) *list.Element {
    return tokenIt
}
