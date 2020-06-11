package statement

import (
    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    block2 "tklibs/script/compiler/ast/statement/block"
    _case "tklibs/script/compiler/ast/statement/case"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type SwitchStatementParserComponent struct {
    script.ComponentType
}

func NewSwitchStatementParserComponent(owner interface{}) *SwitchStatementParserComponent {
    return &SwitchStatementParserComponent{script.MakeComponentType(owner)}
}

func (s *SwitchStatementParserComponent) ParseSwitch(switchStatement interface{}, tokenIt *token.Iterator) *token.Iterator {
    ss := switchStatement.(statement.Switch)
    if tokenIt == nil {
        panic("")
    }

    t := tokenIt.Value().(token.Token)
    if t.GetType() != token.TokenTypeLPAREN {
        panic("switch exception (")
    }

    target, tokenIt := s.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
    ss.SetTargetValue(target)

    t = tokenIt.Value().(token.Token)
    if t.GetType() != token.TokenTypeLBRACE {
        panic("switch exception {")
    }

    tokenIt = tokenIt.Next()

    for tokenIt != nil {
        t = tokenIt.Value().(token.Token)
        switch t.GetType() {
        case token.TokenTypeRBRACE:
            return tokenIt.Next()
        case token.TokenTypeIDENT:
            switch t.GetValue() {
            case "case":
                val, next := s.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt.Next())
                b := &struct {
                    *block2.Component
                }{}

                t = next.Value().(token.Token)
                if t.GetType() != token.TokenTypeCOLON {
                    panic("switch exception :")
                }

                b.Component = block2.NewBlock(b)
                tokenIt = s.GetOwner().(parser.BlockParser).ParseBlock(b, next.Next())
                c := &struct {
                    *_case.Component
                }{}
                c.Component = _case.NewCase(c)
                c.SetValue(val)
                c.SetBlock(b)
                ss.GetCaseList().PushBack(c)
            case "default":
                tokenIt = tokenIt.Next()
                t = tokenIt.Value().(token.Token)
                if t.GetType() != token.TokenTypeCOLON {
                    panic("switch exception :")
                }

                tokenIt = tokenIt.Next()

                b := &struct {
                    *block2.Component
                }{}
                b.Component = block2.NewBlock(b)
                tokenIt = s.GetOwner().(parser.BlockParser).ParseBlock(b, tokenIt)
                ss.SetDefaultCase(b)
            default:
                panic("switch exception 'case' or 'default' ")
            }
        }
    }

    panic("switch invalid statement")
}
