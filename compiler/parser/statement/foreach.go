package statement

import (
    "tklibs/script"
    "tklibs/script/compiler/token"
)

type ForeachStatementParserComponent struct {
    *script.ComponentType
}

func (parser *ForeachStatementParserComponent) ParseForeach(f interface{}, tokenIt *token.Iterator) *token.Iterator {
    return tokenIt
}
