package statement

import (
    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type ForStatementParserComponent struct {
    script.ComponentType
}

func NewForStatementParserComponent(owner interface{}) *ForStatementParserComponent {
    return &ForStatementParserComponent{script.MakeComponentType(owner)}
}

func (impl *ForStatementParserComponent) ParseFor(f interface{}, tokenIt *token.Iterator) *token.Iterator {
    if tokenIt == nil {
        return nil
    }

    if tokenIt.Value().(token.Token).GetType() != token.TokenTypeLPAREN {
        panic("for statement expecting (")
    }

    fs := f.(statement.For)
    blockParser := impl.GetOwner().(parser.BlockParser)

    init, next := blockParser.ParseStatement(tokenIt.Next())
    fs.SetInit(init)
    tokenIt = next

    condition, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
    fs.SetCondition(condition)
    tokenIt = next

    step, next := blockParser.ParseStatement(tokenIt)
    fs.SetStep(step)
    tokenIt = next

    switch tokenIt.Value().(token.Token).GetType() {
    case token.TokenTypeLBRACE:
        body := &struct {
            *block.Component
        }{}
        body.Component = block.NewBlock(body)
        tokenIt = blockParser.ParseBlock(body, tokenIt.Next()).Next()
        fs.SetBody(body)
    default:
        body, next := blockParser.ParseStatement(tokenIt)
        fs.SetBody(body)
        tokenIt = next
    }

    return tokenIt
}
