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

	obj.ScriptSet("call", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return this.(script.Function).Invoke(context, script.Null)
		}

		return this.(script.Function).Invoke(context, args[0], args[1:]...)
	}, impl.GetOwner()))

	obj.ScriptSet("bind", function.NativeFunctionToValue(func(context interface{}, this interface{}, args ...interface{}) interface{} {
		if len(args) < 1 {
			return this
		}

		_func := this.(script.Function)
		if !_func.IsScriptFunction() {
			return this
		}

		f := &struct {
			*function.Component
		}{}

		if _func.IsScriptFunction() {
			f.Component = function.NewScriptFunction(f, _func.GetScriptRuntimeFunction(), context)
		} else {
			f.Component = function.NewScriptFunction(f, _func.GetNativeRuntimeFunction(), context)
		}

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
