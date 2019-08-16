package expression

import (
    "container/list"
    "strconv"

    "tklibs/script"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/expression/arglist"
    "tklibs/script/compiler/ast/expression/array"
    "tklibs/script/compiler/ast/expression/binary"
    "tklibs/script/compiler/ast/expression/call"
    "tklibs/script/compiler/ast/expression/const"
    function2 "tklibs/script/compiler/ast/expression/function"
    _map "tklibs/script/compiler/ast/expression/map"
    "tklibs/script/compiler/ast/expression/member"
    "tklibs/script/compiler/ast/expression/unary"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type tokenTypeLevel struct {
    tokenType token.TokenType
    level     token.TokenType
}

func getMaxOpLevel(opList *list.List) token.TokenType {
    level := token.TokenTypeUnknown

    for it := opList.Front(); it != nil; it = it.Next() {
        op := it.Value.(*tokenTypeLevel)
        if op.level > level {
            level = op.level
        }
    }

    return level
}

func makeExpression(expressionList *list.List, op *tokenTypeLevel, e interface{}) {
    cur := expressionList.Back().Value
    expressionList.Remove(expressionList.Back())

    switch op.tokenType {
    case token.TokenTypePERIOD:
        m := &struct {
            *member.Component
        }{}
        m.Component = member.NewMember(m, e, cur)
        expressionList.PushBack(m)
    case token.TokenTypeLPAREN:
        c := &struct {
            *call.Component
        }{}
        c.Component = call.NewCall(c, e, cur)
        expressionList.PushBack(c)
    case token.TokenTypeSUB,
        token.TokenTypeLNOT:
        if e == nil {
            u := &struct {
                *unary.Component
            }{}
            u.Component = unary.NewUnary(u, cur, op.tokenType)
            expressionList.PushBack(u)
            break
        }
        fallthrough
    default:
        b := &struct {
            *binary.Component
        }{}
        b.Component = binary.NewBinary(b, e, cur, op.tokenType)
        expressionList.PushBack(b)
    }
}

func processExpression(opList, expressionList *list.List) {
    if opList.Len() > 0 && expressionList.Len() > 0 {
        e := expressionList.Back().Value
        expressionList.Remove(expressionList.Back())

        op := opList.Back().Value.(*tokenTypeLevel)
        opList.Remove(opList.Back())

        maxOpLevel := getMaxOpLevel(opList)

        if op.level < maxOpLevel {
            processExpression(opList, expressionList)
            opList.PushBack(op)
            expressionList.PushBack(e)
        } else {
            makeExpression(expressionList, op, e)
        }
    }
}

func getExpression(opList, expressionList *list.List) interface{} {
    for opList.Len() > 0 {
        e := expressionList.Back().Value
        expressionList.Remove(expressionList.Back())

        op := opList.Back().Value.(*tokenTypeLevel)
        opList.Remove(opList.Back())

        maxOpLevel := getMaxOpLevel(opList)

        if op.level < maxOpLevel && opList.Len() > 0 {
            processExpression(opList, expressionList)
            opList.PushBack(op)
            expressionList.PushBack(e)
            continue
        }

        makeExpression(expressionList, op, e)
    }

    if expressionList.Len() == 0 {
        return nil
    }

    return expressionList.Front().Value
}

type ExpressionParserComponent struct {
    script.ComponentType
}

func NewExpressionParser(owner interface{}) *ExpressionParserComponent {
    return &ExpressionParserComponent{script.MakeComponentType(owner)}
}

func (p *ExpressionParserComponent) ParseExpression(tokenIt *list.Element) (interface{}, *list.Element) {
    //todo check null for tokenIt

    expressionList := list.New()
    opList := list.New()

parseLoop:
    for {
        if tokenIt == nil {
            break
        }

        var currentExpression, currentOp interface{}

        t := tokenIt.Value.(token.Token)

        tokenType := t.GetType()
        switch tokenType {
        case token.TokenTypeCOMMA: // ,
            break parseLoop
        case token.TokenTypeLBRACE: // {
            if opList.Len() < expressionList.Len() {
                break parseLoop
            }

            m := &struct {
                *_map.Component
            }{}
            m.Component = _map.NewMap(m)
            currentExpression, tokenIt = m, p.GetOwner().(parser.MapParser).ParseMap(m, tokenIt.Next())
        case token.TokenTypeLBRACK: // [
            if opList.Len() < expressionList.Len() {
                e := expressionList.Front().Value
                expressionList.Remove(expressionList.Front())
                field, next := p.ParseExpression(tokenIt.Next())
                m := &struct {
                    *member.Component
                }{}
                m.Component = member.NewMember(m, e, field)
                currentExpression, tokenIt = m, next
            } else {
                a := &struct {
                    *array.ArrayExpressionComponent
                }{}
                al := &struct {
                    *arglist.Component
                }{}
                al.Component = arglist.NewArgList(al)
                a.ArrayExpressionComponent = array.NewArrayExpression(a, al)
                currentExpression, tokenIt = a, p.GetOwner().(parser.ArgListParser).ParseArgList(al, tokenIt.Next())
            }
        case token.TokenTypeRBRACK, token.TokenTypeRPAREN, token.TokenTypeSEMICOLON:
            tokenIt = tokenIt.Next()
            break parseLoop
        case token.TokenTypeIDENT:
            if opList.Len() < expressionList.Len() {
                break parseLoop
            }
            value := t.GetValue()
            switch value {
            case "new":
                exp, next := p.ParseExpression(tokenIt.Next())
                makeNewCall(exp)
                currentExpression, tokenIt = exp, next
            case "function", "func", "#":
                f := &struct {
                    *function2.Component
                }{}
                f.Component = function2.NewFunction(f)
                currentExpression, tokenIt = f, p.GetOwner().(parser.FunctionParser).ParseFunction(f, tokenIt.Next())
            case "yes", "no", "true", "false":
                c := &struct {
                    *_const.Component
                }{}

                value := t.GetValue()

                if value == "yes" || value == "true" {
                    c.Component = _const.NewConst(c, script.Bool(true))
                } else {
                    c.Component = _const.NewConst(c, script.Bool(false))
                }
                currentExpression, tokenIt = c, tokenIt.Next()
            default:
                if opList.Len() < expressionList.Len() {
                    break parseLoop
                }

                if opList.Len() > 0 {
                    switch opList.Front().Value.(*tokenTypeLevel).tokenType {
                    case token.TokenTypePERIOD:
                        m := &struct {
                            *member.Component
                        }{}
                        m.Component = member.NewMember(m, expressionList.Front().Value, t.GetValue())
                        opList.Remove(opList.Front())
                        expressionList.Remove(expressionList.Front())

                        currentExpression, tokenIt = m, tokenIt.Next()
                    default:
                        m := &struct {
                            *member.Component
                        }{}
                        m.Component = member.NewMember(m, nil, t.GetValue())
                        currentExpression, tokenIt = m, tokenIt.Next()
                    }
                } else {
                    m := &struct {
                        *member.Component
                    }{}
                    m.Component = member.NewMember(m, nil, t.GetValue())
                    currentExpression, tokenIt = m, tokenIt.Next()
                }
            }
        case token.TokenTypeINT:
            value, _ := strconv.ParseInt(t.GetValue(), 10, 64)
            c := &struct {
                *_const.Component
            }{}
            c.Component = _const.NewConst(c, script.Int(value))
            currentExpression, tokenIt = c, tokenIt.Next()
        case token.TokenTypeFLOAT:
            value, _ := strconv.ParseFloat(t.GetValue(), 64)
            c := &struct {
                *_const.Component
            }{}

            c.Component = _const.NewConst(c, script.Float(value))
            currentExpression, tokenIt = c, tokenIt.Next()
        case token.TokenTypeSTRING:
            c := &struct {
                *_const.Component
            }{}
            c.Component = _const.NewConst(c, script.String(t.GetValue()))
            currentExpression, tokenIt = c, tokenIt.Next()
        case token.TokenTypeINC, token.TokenTypeDEC:
            ce := expressionList.Back().Value
            expressionList.Remove(expressionList.Back())

            u := &struct {
                *unary.Component
            }{}
            u.Component = unary.NewUnary(u, ce, t.GetType())

            currentExpression, tokenIt = u, tokenIt.Next()
        case
            token.TokenTypeAND,
            token.TokenTypeANDASSIGN,
            token.TokenTypeOR,
            token.TokenTypeORASSIGN,
            token.TokenTypeXOR,
            token.TokenTypeXORASSIGN,
            token.TokenTypeADD,
            token.TokenTypeADDASSIGN,
            token.TokenTypeSUB,
            token.TokenTypeSUBASSIGN,
            token.TokenTypeMUL,
            token.TokenTypeMULASSIGN,
            token.TokenTypeDIV,
            token.TokenTypeDIVASSIGN,
            token.TokenTypeSHR,
            token.TokenTypeSHRASSIGN,
            token.TokenTypeSHL,
            token.TokenTypeSHLASSIGN,
            token.TokenTypeLOR,
            token.TokenTypeLAND,
            token.TokenTypeLNOT,
            token.TokenTypeEQL,
            token.TokenTypeNEQ,
            token.TokenTypeGREATER,
            token.TokenTypeGEQ,
            token.TokenTypeLESS,
            token.TokenTypeLEQ,
            token.TokenTypeREM,
            token.TokenTypeREMASSIGN,
            token.TokenTypeASSIGN:
            if opList.Len() == expressionList.Len() {
                switch tokenType {
                case token.TokenTypeSUB,
                    token.TokenTypeLNOT:
                    expressionList.PushFront(nil)
                    currentOp, tokenIt = &tokenTypeLevel{tokenType, token.TokenTypeCOLON}, tokenIt.Next()
                default:
                    currentOp, tokenIt = &tokenTypeLevel{tokenType, tokenType}, tokenIt.Next()
                }
            } else {
                currentOp, tokenIt = &tokenTypeLevel{tokenType, tokenType}, tokenIt.Next()
            }
        case token.TokenTypePERIOD: // .
            currentOp, tokenIt = &tokenTypeLevel{token.TokenTypePERIOD, token.TokenTypePERIOD}, tokenIt.Next()
        case token.TokenTypeLPAREN: // (
            if opList.Len() < expressionList.Len() {
                // it's a call expression
                a := &struct {
                    *arglist.Component
                }{}
                a.Component = arglist.NewArgList(a)
                argList, next := a, p.GetOwner().(parser.ArgListParser).ParseArgList(a, tokenIt.Next())

                c := &struct {
                    *call.Component
                }{}
                c.Component = call.NewCall(c, expressionList.Front().Value, argList)
                expressionList.Remove(expressionList.Front())

                currentExpression, tokenIt = c, next
            } else {
                // check for lambda function
                it := tokenIt
                for ; it != nil; it = it.Next() { // scan token list until ')'
                    if it.Value.(token.Token).GetType() == token.TokenTypeRPAREN {
                        break
                    }
                }
                // check the next symbol '=>'
                if it != nil && it.Next() != nil && it.Next().Value.(token.Token).GetType() == token.TokenTypeLAMBDA {
                    f := &struct {
                        *function2.Component
                    }{}
                    f.Component = function2.NewFunction(f)
                    currentExpression, tokenIt = f, p.GetOwner().(parser.FunctionParser).ParseFunction(f, tokenIt)
                } else {
                    currentExpression, tokenIt = p.ParseExpression(tokenIt.Next())
                }
            }
        default:
            break parseLoop
        }

        if currentOp != nil {
            opList.PushFront(currentOp)
        }

        if currentExpression != nil {
            expressionList.PushFront(currentExpression)
        }
    }

    return getExpression(opList, expressionList), tokenIt
}

func makeNewCall(exp interface{}) bool {
    switch c := exp.(type) {
    case expression.Call:
        if !makeNewCall(c.GetExpression()) {
            c.SetNew(true)
            return true
        }
    case expression.Member:
        if c, ok := c.GetLeft().(expression.Call); ok {
            return makeNewCall(c)
        }
    case expression.Binary:
        if c, ok := c.GetLeft().(expression.Call); ok {
            return makeNewCall(c)
        }
    }

    return false
}
