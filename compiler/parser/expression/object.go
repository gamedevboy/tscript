package expression

import (
	"tklibs/script"
	"tklibs/script/compiler/ast/expression"
	"tklibs/script/compiler/parser"
	"tklibs/script/compiler/token"
)

type ObjectExpressionParserComponent struct {
	script.ComponentType
}

func NewObjectExpressionParser(owner interface{}) *ObjectExpressionParserComponent {
	return &ObjectExpressionParserComponent{script.MakeComponentType(owner)}
}

func (p *ObjectExpressionParserComponent) ParseObject(m interface{}, tokenIt *token.Iterator) *token.Iterator {
	if tokenIt == nil {
		return nil
	}

	object := m.(expression.Object)

	for {
		t := tokenIt.Value().(token.Token)

		switch t.GetType() {
		case token.TokenTypeCOMMA: // skip ,
			tokenIt = tokenIt.Next()
			continue
		case token.TokenTypeRBRACE:
			return tokenIt.Next()
		case token.TokenTypeIDENT:
			varName := t.GetValue() // get var name

			tokenIt = tokenIt.Next() // skip :
			// TODO check if the token is = or :
			if tokenIt == nil {
				return nil
			}

			t = tokenIt.Value().(token.Token)
			if t.GetType() != token.TokenTypeCOLON {
				return tokenIt
			}

			e, next := p.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt.Next()) // parse expression
			*object.GetKeyValueMap() = append(*object.GetKeyValueMap(), expression.ObjectEntry{
				varName,
				e,
			})
			tokenIt = next
		default:
			panic("")
		}
	}
}
