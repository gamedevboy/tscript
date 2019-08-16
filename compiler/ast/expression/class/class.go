package class

import (
    "container/list"

    "tklibs/script"
    "tklibs/script/compiler"
    "tklibs/script/compiler/ast/expression"
    "tklibs/script/compiler/ast/expression/arglist"
    "tklibs/script/compiler/ast/expression/binary"
    "tklibs/script/compiler/ast/expression/call"
    "tklibs/script/compiler/ast/expression/function"
    "tklibs/script/compiler/ast/expression/member"
    "tklibs/script/compiler/ast/statement"
    "tklibs/script/compiler/ast/statement/block"
    "tklibs/script/compiler/ast/statement/decl"
    expression2 "tklibs/script/compiler/ast/statement/expression"
    _return "tklibs/script/compiler/ast/statement/return"
    "tklibs/script/compiler/token"
)

type callComponent = call.Component

type Component struct {
    script.ComponentType
    *callComponent
    function       interface{}
    methods        list.List
    parent         interface{}
    name           string
    hasConstructor bool
}

func (impl *Component) SetName(name string) {
    impl.name = name
}

func (impl *Component) GetName() string {
    return impl.name
}

func (impl *Component) SetParent(parent interface{}) {
    impl.parent = parent
}

func (impl *Component) GetParent() interface{} {
    return impl.parent
}

func (impl *Component) GetMethods() *list.List {
    return &impl.methods
}

func (impl *Component) AddMethod(f interface{}) {
    impl.methods.PushBack(f)
    _f := f.(expression.Function)
    b := impl.function.(expression.Function).GetBlock().(statement.Block)

    if _f.GetName() == "constructor" {
        if impl.hasConstructor {
            panic("only one constructor allowed in class")
        }

        d := &struct{ *decl.Component }{}
        d.Component = decl.NewDecl(d)

        d.SetName(_f.GetName())
        d.SetExpression(f)
        b.GetStatementList().PushFront(d)
        impl.hasConstructor = true
    } else {
        c := &struct {
            *member.Component
        }{}
        c.Component = member.NewMember(c, nil, "constructor")

        m := &struct {
            *member.Component
        }{}
        m.Component = member.NewMember(m, c, _f.GetName())

        bin := &struct {
            *binary.Component
        }{}
        bin.Component = binary.NewBinary(bin, m, f, token.TokenTypeASSIGN)

        es := &struct {
            *expression2.Component
        }{}
        es.Component = expression2.NewExpressionStatement(es)
        es.SetExpression(bin)

        b.GetStatementList().PushBack(es)
    }
}

func (impl *Component) FinishMethod() {
    _block := impl.function.(expression.Function).GetBlock().(statement.Block)

    if !impl.hasConstructor {
        f := &struct {
            *function.Component
        }{}
        f.Component = function.NewFunction(f)
        f.SetName("constructor")

        bl := &struct {
            *block.Component
        }{}
        bl.Component = block.NewBlock(bl)

        f.SetBlock(bl)

        d := &struct{ *decl.Component }{}
        d.Component = decl.NewDecl(d)
        d.SetExpression(f)
        d.SetName("constructor")
        _block.GetStatementList().PushFront(d)
    }

    r := &struct {
        *_return.Component
    }{}
    r.Component = _return.NewReturn(r)

    rv := &struct {
        *member.Component
    }{}
    rv.Component = member.NewMember(rv, nil, "constructor")

    r.SetExpression(rv)

    if impl.GetParent() != nil {
        p := &struct {
            *member.Component
        }{}

        p.Component = member.NewMember(p, rv, "prototype")

        b := &struct {
            *binary.Component
        }{}
        b.Component = binary.NewBinary(b, p, impl.GetParent(), token.TokenTypeASSIGN)

        es := &struct {
            *expression2.Component
        }{}
        es.Component = expression2.NewExpressionStatement(es)
        es.SetExpression(b)

        _block.GetStatementList().PushBack(es)
    }

    _block.GetStatementList().PushBack(r)
}

func NewComponent(owner interface{}) *Component {
    f := &struct{ *function.Component }{}
    f.Component = function.NewFunction(f)

    b := &struct{ *block.Component }{}
    b.Component = block.NewBlock(b)

    f.SetBlock(b)

    al := &struct{ *arglist.Component }{}
    al.Component = arglist.NewArgList(al)

    ret := &Component{
        ComponentType: script.MakeComponentType(owner),
        function:      f,
        callComponent: call.NewCall(owner, f, al),
    }

    return ret
}

func (impl *Component) Compile(f interface{}, r *compiler.Operand) *compiler.Operand {
    _func := f.(compiler.Function)
    if r == nil {
        r = compiler.NewRegisterOperand(_func.AllocRegister(""))
    }

    return impl.callComponent.Compile(f, r)
}
