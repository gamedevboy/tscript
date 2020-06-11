package statement

import (
    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type WhileStatementParserComponent struct {
    script.ComponentType
}

func (impl *WhileStatementParserComponent) ParseWhile(w interface{}, tokenIt *token.Iterator) *token.Iterator {
    if tokenIt == nil {
        return tokenIt
    }

    if tokenIt.Value().(token.Token).GetType() != token.TokenTypeLPAREN {
        panic("while expecting ( ")
    }

    ws := w.(statement.While)

    e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
    ws.SetCondition(e)
    tokenIt = next

    blockParser := impl.GetOwner().(parser.BlockParser)
    switch tokenIt.Value().(token.Token).GetType() {
    case token.TokenTypeLBRACE:
        body := &struct {
            *block.Component
        }{}
        body.Component = block.NewBlock(body)
        tokenIt = blockParser.ParseBlock(body, tokenIt.Next()).Next()
        ws.SetBody(body)
    default:
        body, next := blockParser.ParseStatement(tokenIt)
        ws.SetBody(body)
        tokenIt = next
    }

    return tokenIt
}

func NewWhileStatementParserComponent(owner interface{}) *WhileStatementParserComponent {
    return &WhileStatementParserComponent{script.MakeComponentType(owner)}
}
