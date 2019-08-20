package statement

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler/ast/expression/class"
    "tklibs/script/compiler/ast/expression/function"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    _break "tklibs/script/compiler/ast/statement/break"
    _continue "tklibs/script/compiler/ast/statement/continue"
    "tklibs/script/compiler/ast/statement/decl"
    "tklibs/script/compiler/ast/statement/expression"
    "tklibs/script/compiler/ast/statement/for"
    "tklibs/script/compiler/ast/statement/if"
    "tklibs/script/compiler/ast/statement/return"
    _switch "tklibs/script/compiler/ast/statement/switch"
    "tklibs/script/compiler/ast/statement/while"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type BlockStatementParserComponent struct {
    script.ComponentType
}

func (impl *BlockStatementParserComponent) ParseStatement(tokenIt *list.Element) (interface{}, *list.Element) {
    if tokenIt == nil {
        return nil, nil
    }

    t := tokenIt.Value.(token.Token)
    switch t.GetType() {
    case token.TokenTypeIDENT:
        switch t.GetValue() {
        case "var":
            ds := &struct {
                *decl.Component
            }{}
            ds.Component = decl.NewDecl(ds)
            ds.SetLine(t.GetLine())
            ds.SetFilePath(t.GetFilePath())
            return ds, impl.GetOwner().(parser.DeclParser).ParseDecl(ds, tokenIt.Next())
        case "for":
            f := &struct {
                *_for.Component
            }{}
            f.Component = _for.NewForStatementComponent(f)
            f.SetLine(t.GetLine())
            f.SetFilePath(t.GetFilePath())
            return f, impl.GetOwner().(parser.ForParser).ParseFor(f, tokenIt.Next())
        case "foreach":
            panic(fmt.Errorf("Not supported yet ! "))
        case "while":
            ws := &struct {
                *while.Component
            }{}
            ws.Component = while.NewWhileStatementComponent(ws)
            ws.SetLine(t.GetLine())
            ws.SetFilePath(t.GetFilePath())
            return ws, impl.GetOwner().(parser.WhileParser).ParseWhile(ws, tokenIt.Next())
        case "if":
            is := &struct {
                *_if.Component
            }{}
            is.Component = _if.NewIfStatementComponent(is)
            is.SetLine(t.GetLine())
            is.SetFilePath(t.GetFilePath())
            return is, impl.GetOwner().(parser.IfParser).ParseIf(is, tokenIt.Next())
        case "break":
            b := &struct {
                *_break.Component
            }{}
            b.Component = _break.NewBreak(b)
            return b, tokenIt.Next()
        case "continue":
            c := &struct {
                *_continue.Component
            }{}
            c.Component = _continue.NewContinue(c)
            return c, tokenIt.Next()
        case "switch":
            s := &struct {
                *_switch.Component
            }{}
            s.Component = _switch.NewSwitch(s)
            return s, impl.GetOwner().(parser.SwitchParser).ParseSwitch(s, tokenIt.Next())
        case "return":
            rs := &struct {
                *_return.Component
            }{}
            rs.Component = _return.NewReturn(rs)
            rs.SetLine(t.GetLine())
            rs.SetFilePath(t.GetFilePath())

            next := tokenIt.Next()
            if next != nil && next.Value.(token.Token).GetLine() == t.GetLine() {
                e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(next)
                rs.SetExpression(e)
                return rs, next
            }

            return rs, next
        case "class":
            ds := &struct {
                *decl.Component
            }{}
            ds.Component = decl.NewDecl(ds)
            ds.SetLine(t.GetLine())
            ds.SetFilePath(t.GetFilePath())
            c := &struct {
                *class.Component
            }{}
            c.Component = class.NewComponent(c)
            next := impl.GetOwner().(parser.ClassParser).ParseClass(c, tokenIt.Next())
            ds.SetName(c.GetName())
            ds.SetExpression(c)
            return ds, next
        case "function", "func", "#":
            ds := &struct {
                *decl.Component
            }{}
            ds.Component = decl.NewDecl(ds)
            ds.SetLine(t.GetLine())
            ds.SetFilePath(t.GetFilePath())
            f := &struct {
                *function.Component
            }{}
            f.Component = function.NewFunction(f)
            next := impl.GetOwner().(parser.FunctionParser).ParseFunction(f, tokenIt.Next())
            ds.SetName(f.Component.GetName())
            ds.SetExpression(f)
            return ds, next
        default:
            e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
            es := &struct {
                *expression.Component
            }{}
            es.Component = expression.NewExpressionStatement(es)
            es.SetLine(t.GetLine())
            es.SetFilePath(t.GetFilePath())
            es.SetExpression(e)
            return es, next
        }
    case token.TokenTypeLBRACE:
        _bs := &struct {
            *block.Component
        }{}
        _bs.Component = block.NewBlock(_bs)
        return _bs, impl.ParseBlock(_bs, tokenIt.Next()).Next()
    default:
        e, next := impl.GetOwner().(parser.ExpressionParser).ParseExpression(tokenIt)
        es := &struct {
            *expression.Component
        }{}
        es.Component = expression.NewExpressionStatement(es)
        es.SetExpression(e)
        return es, next
    }

    return nil, nil
}

func (impl *BlockStatementParserComponent) ParseBlock(bs interface{}, tokenIt *list.Element) *list.Element {
    if tokenIt == nil {
        return nil
    }
    blockStatement := bs.(statement.Block)

    it := tokenIt
    for it != nil {
        t := it.Value.(token.Token)
        switch t.GetType() {
        case token.TokenTypeRBRACE:
            return it
        case token.TokenTypeIDENT:
            if t.GetValue() == "case" {
                return it
            }
            if t.GetValue() == "default" {
                return it
            }
            fallthrough
        default:
            es, next := impl.ParseStatement(it)
            if es != nil {
                blockStatement.GetStatementList().PushBack(es)
            }

            it = next
        }
    }

    return it
}

func NewBlockStatementParser(owner interface{}) *BlockStatementParserComponent {
    return &BlockStatementParserComponent{script.MakeComponentType(owner)}
}
