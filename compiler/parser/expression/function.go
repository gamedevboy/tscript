package expression

import (
    "container/list"
    "fmt"

    "tklibs/script"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/expression/arglist"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/debug"
    "tklibs/script/compiler/parser"
    "tklibs/script/compiler/token"
)

type FunctionExpressionParserComponent struct {
    script.ComponentType
}

func NewFunctionExpressionParser(owner interface{}) *FunctionExpressionParserComponent {
    return &FunctionExpressionParserComponent{script.MakeComponentType(owner)}
}

func (impl *FunctionExpressionParserComponent) ParseFunction(f interface{}, tokenIt *list.Element) *list.Element {
    if tokenIt == nil {
        return nil
    }
    _func := f.(expression.Function)
    t := tokenIt.Value.(token.Token)
    debugInfo := f.(debug.Info)
    debugInfo.SetLine(t.GetLine())
    debugInfo.SetFilePath(t.GetFilePath())
    for i := 0; ; i++ {
        if tokenIt == nil {
            return nil
        }

        t := tokenIt.Value.(token.Token)

        switch t.GetType() {
        case token.TokenTypeLAMBDA:

        case token.TokenTypeLPAREN:
            a := &struct {
                *arglist.Component
            }{}
            a.Component = arglist.NewArgList(a)

            fe := _func
            fe.SetArgList(a)

            tokenIt = impl.GetOwner().(parser.ArgListParser).ParseArgList(a, tokenIt.Next())

            t = tokenIt.Value.(token.Token)

            if t.GetType() == token.TokenTypeLAMBDA {
                fe.SetCaptureThis(true)
                tokenIt = tokenIt.Next()
                t = tokenIt.Value.(token.Token)
            }

            if t.GetType() != token.TokenTypeLBRACE {
                panic("")
            }
            b := &struct {
                *block.Component
            }{}
            b.Component = block.NewBlock(b)
            fe.SetBlock(b)
            return impl.GetOwner().(parser.BlockParser).ParseBlock(b, tokenIt.Next()).Next()
        default:
            switch i {
            case 0:
                _func.SetName(t.GetValue())
            default:
                panic(fmt.Errorf("%v: %v ==> function excepting (", t.GetFilePath(), t.GetLine()))
            }
        }

        tokenIt = tokenIt.Next()
    }
}
