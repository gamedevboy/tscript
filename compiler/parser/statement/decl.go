package statement

import (
    "tklibs/script"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type DeclStatementParserComponent struct {
    script.ComponentType
}

func (impl *DeclStatementParserComponent) ParseDecl(ds interface{}, tokenIt *token.Iterator) *token.Iterator {
    if tokenIt == nil {
        panic("wrong decl, excepting }")
    }

    varName := tokenIt
    ds.(statement.Decl).SetName(varName.Value().(token.Token).GetValue())

    op := varName.Next()

    if op == nil {
        return nil
    }

    opToken := op.Value().(token.Token)
    switch opToken.GetType() {
    case token.TokenTypeASSIGN:
        e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(op.Next())
        ds.(statement.Decl).SetExpression(e)
        return next
    case token.TokenTypeSEMICOLON:
        return op.Next()
    default:
        return op
    }
}

func NewDeclStatementParser(owner interface{}) *DeclStatementParserComponent {
    return &DeclStatementParserComponent{script.MakeComponentType(owner)}
}
