package expression

import (
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/expression/function"
    "tklibs/script/compiler/ast/expression/member"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type ClassExpressionParserComponent struct {
    script.ComponentType
}

func NewClassExpressionParser(owner interface{}) *ClassExpressionParserComponent {
    return &ClassExpressionParserComponent{script.MakeComponentType(owner)}
}

func (impl *ClassExpressionParserComponent) parse(c interface{}, tokenIt *token.Iterator) *token.Iterator {
    e := c.(expression.Class)
    methods := e.GetMethods()

    for tokenIt != nil {
        t := tokenIt.Value().(token.Token)
        switch t.GetType() {
        case token.TokenTypeRBRACE:
            e.FinishMethod()
            return tokenIt.Next()
        default:
            f := &struct {
                *function.Component
            }{}
            f.Component = function.NewFunction(f)
            methods.PushBack(f)
            tokenIt = impl.GetOwner().(parser.FunctionParser).ParseFunction(f, tokenIt)
            e.AddMethod(f)
        }
    }

    e.FinishMethod()

    return nil
}

func (impl *ClassExpressionParserComponent) ParseClass(c interface{}, tokenIt *token.Iterator) *token.Iterator {
    if tokenIt == nil {
        return tokenIt
    }

    t := tokenIt.Value().(token.Token)
    className := t.GetValue()
    c.(expression.Class).SetName(className)

    tokenIt = tokenIt.Next()
    if tokenIt == nil {
        return nil
    }

    t = tokenIt.Value().(token.Token)
    switch t.GetType() {
    case token.TokenTypeCOLON:
        tokenIt = tokenIt.Next()
        t = tokenIt.Value().(token.Token)

        parentName := t.GetValue()
        m := &struct {
            *member.Component
        }{}
        m.Component = member.NewMember(m, nil, parentName, false)
        c.(expression.Class).SetParent(m)

        tokenIt = tokenIt.Next()

        t = tokenIt.Value().(token.Token)
        if t.GetType() != token.TokenTypeLBRACE {
            panic("excepting {")
        }

        return impl.parse(c, tokenIt.Next())
    case token.TokenTypeLBRACE:
        return impl.parse(c, tokenIt.Next())
    default:
        panic(fmt.Errorf("excepting : "))
    }
}
