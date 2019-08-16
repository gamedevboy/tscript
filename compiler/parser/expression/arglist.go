package expression

import (
    "container/list"
    "fmt"

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

func (impl *ArgListExpressionParserComponent) ParseArgList(a interface{}, tokenIt *list.Element) *list.Element {
    for {
        if tokenIt == nil {
            return nil
        }

        e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
        if e != nil {
            a.(expression.ArgList).GetExpressionList().PushBack(e)
        } else {
            return next
        }

        tokenIt = next

        if tokenIt != nil {
            prev := tokenIt.Prev()
            if prev.Value.(token.Token).GetType() == token.TokenTypeRPAREN {
                tokenIt = prev
            }

            t := tokenIt.Value.(token.Token)

            switch t.GetType() {
            case token.TokenTypeCOMMA:
                tokenIt = tokenIt.Next()
                continue
            case token.TokenTypeRPAREN:
                if _, ok := e.(expression.Call); ok {
                    tokenIt = tokenIt.Next()

                    if tokenIt != nil {
                        switch tokenIt.Value.(token.Token).GetType() {
                        case token.TokenTypeCOMMA:
                            tokenIt = tokenIt.Next()
                        default:
                            return tokenIt
                        }
                    }
                } else {
                    return tokenIt.Next()
                }
            default:
                panic(fmt.Errorf("unexcept token: %v", t.GetValue()))
            }
        }
    }
}
