package parser

import (
    "tklibs/script/compiler/parser/expression"
    "tklibs/script/compiler/parser/statement"
)

type parser struct {
    *expression.FunctionExpressionParserComponent
    *expression.ExpressionParserComponent
    *expression.ArgListExpressionParserComponent
    *expression.MapExpressionParserComponent
    *expression.ClassExpressionParserComponent
    *statement.BlockStatementParserComponent
    *statement.DeclStatementParserComponent
    *statement.IfStatementParserComponent
    *statement.ForStatementParserComponent
    *statement.WhileStatementParserComponent
    *statement.SwitchStatementParserComponent
}

func NewParser() interface{} {
    p := &parser{}

    p.ClassExpressionParserComponent = expression.NewClassExpressionParser(p)
    p.FunctionExpressionParserComponent = expression.NewFunctionExpressionParser(p)
    p.ExpressionParserComponent = expression.NewExpressionParser(p)
    p.ArgListExpressionParserComponent = expression.NewArgListExpressionParser(p)
    p.MapExpressionParserComponent = expression.NewMapExpressionParser(p)
    p.BlockStatementParserComponent = statement.NewBlockStatementParser(p)
    p.DeclStatementParserComponent = statement.NewDeclStatementParser(p)
    p.IfStatementParserComponent = statement.NewIfStatementParserComponent(p)
    p.ForStatementParserComponent = statement.NewForStatementParserComponent(p)
    p.WhileStatementParserComponent = statement.NewWhileStatementParserComponent(p)
    p.SwitchStatementParserComponent = statement.NewSwitchStatementParserComponent(p)

    return p
}
