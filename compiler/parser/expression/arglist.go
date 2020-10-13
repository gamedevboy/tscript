package expression

import (
    "tklibs/script"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type ArgListExpressionParserComponent struct {
    script.ComponentType
}

func NewArgListExpressionParser(owner interface{}) *ArgListExpressionParserComponent {
    return &ArgListExpressionParserComponent{script.MakeComponentType(owner)}
}

func (impl *ArgListExpressionParserComponent) ParseArgList(a interface{}, tokenIt *token.Iterator) *token.Iterator {
    for {
        if tokenIt == nil {
            return nil
        }

        if tokenIt.Value().(token.Token).GetType() == token.TokenTypeRPAREN {
            return tokenIt.Next()
        }

        e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
        if e != nil {
            a.(expression.ArgList).GetExpressionList().PushBack(e)
        } else {
            return next
        }

        tokenIt = next

        if tokenIt != nil {
            t := tokenIt.Value().(token.Token)

            switch t.GetType() {
            case token.TokenTypeCOMMA:
                tokenIt = tokenIt.Next()
                continue
            case token.TokenTypeRPAREN, token.TokenTypeRBRACK:
                return tokenIt.Next()
            default:
                return tokenIt
            }
        }
    }
}
