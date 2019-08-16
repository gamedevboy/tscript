package statement

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type IfStatementParserComponent struct {
    script.ComponentType
}

func NewIfStatementParserComponent(owner interface{}) *IfStatementParserComponent {
    return &IfStatementParserComponent{script.MakeComponentType(owner)}
}

func (impl *IfStatementParserComponent) ParseIf(ifStatement interface{}, tokenIt *list.Element) *list.Element {
    is := ifStatement.(statement.If)
    condition, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
    is.SetCondition(condition)
    tokenIt = next

    switch tokenIt.Value.(token.Token).GetType() {
    case token.TokenTypeLBRACE:
        body := &struct {
            *block.Component
        }{}
        body.Component = block.NewBlock(body)
        tokenIt = impl.GetOwner().(parser.BlockParser).ParseBlock(body, tokenIt.Next()).Next()
        is.SetBody(body)
    default:
        body, next := impl.GetOwner().(parser.BlockParser).ParseStatement(tokenIt)
        is.SetBody(body)
        tokenIt = next
    }

    if tokenIt != nil && tokenIt.Value.(token.Token).GetValue() == "else" {
        tokenIt = tokenIt.Next()

        switch tokenIt.Value.(token.Token).GetType() {
        case token.TokenTypeLBRACE:
            elseBody := &struct {
                *block.Component
            }{}
            elseBody.Component = block.NewBlock(elseBody)
            tokenIt = impl.GetOwner().(parser.BlockParser).ParseBlock(elseBody, tokenIt.Next()).Next()
            is.SetElseBody(elseBody)
        default:
            elseBody, next := impl.GetOwner().(parser.BlockParser).ParseStatement(tokenIt)
            is.SetElseBody(elseBody)
            tokenIt = next
        }
    }

    return tokenIt
}
