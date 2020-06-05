package prototype

import (
	"tklibs/script"
	"tklibs/script/runtime"
	"tklibs/script/type/function"
)

type Function struct {
	script.ComponentType
	prototype interface{}
}

func (impl *Function) GetFunctionPrototype() interface{} {
	return impl.prototype
}

func (impl *Function) InitPrototype() {
	obj := impl.prototype.(script.Object)

	obj.ScriptSet("call", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return this.(script.Function).Invoke(nil)
		}

		return this.(script.Function).Invoke(args[0], args[1:]...)
	}, impl.GetOwner()))

	obj.ScriptSet("bind", function.NativeFunctionToValue(func(this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return this.(script.Function).Invoke(nil)
		}

		_func := this.(script.Function)
		if !_func.IsScriptFunction() {
			return this
		}
		context := _func.GetContext()

		f := &struct {
			*function.Component
		}{}
		f.Component = function.NewScriptFunction(f, _func.GetRuntimeFunction(), context)
		f.SetThis(script.InterfaceToValue(this))
		f.Init()

		return f
	}, impl.GetOwner()))
}

func NewFunctionPrototype(ctx interface{}) *Function {
	ret := &Function{
		ComponentType: script.MakeComponentType(ctx),
	}
	ret.prototype = ctx.(runtime.ScriptContext).NewScriptObject(0)

	return ret
}
