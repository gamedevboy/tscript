package math_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"tklibs/script"
	"tklibs/script/assembly/assembly"
	"tklibs/script/assembly/loader"
	"tklibs/script/runtime/context"
)

var scriptTest = `
function maxInt32() {
	return  math.maxInt32()
}
function toInt(val) {
	return  math.toInt(val)
}
`

var gScriptContext *context.Component

func init() {
	scriptAssembly := &struct {
		*assembly.Component
	}{}

	scriptAssembly.Component = assembly.NewScriptAssembly(scriptAssembly)
	if err := loader.LoadAssemblySourceFromBuffer(scriptAssembly, strings.NewReader(scriptTest)); err != nil {
		panic(err)
	}

	scriptContext := &struct {
		*context.Component
	}{}
	scriptContext.Component = context.NewScriptContext(scriptContext, scriptAssembly, 64)
	gScriptContext = scriptContext.Component
	scriptContext.Run()
}

func checkEnv(t *testing.T, cc *context.Component, fName string, invoker func(script.Function) interface{}) interface{} {
	var tf = cc.ScriptGet(fName).GetFunction()
	if tf == nil {
		t.Errorf("Failed:%s func invalid", fName)
		return nil
	}
	ret := invoker(tf)
	if ret == nil {
		t.Errorf("Failed:%s return nil", fName)
		return nil
	}
	return ret
}

func TestMathMax32(t *testing.T) {
	f := "maxInt32"
	ret := checkEnv(t, gScriptContext, f, func(tf script.Function) interface{} {
		return tf.Invoke(nil)
	})
	if ret != nil && int32(ret.(script.Int)) != math.MaxInt32 {
		t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), math.MaxInt32)
		return
	}
}

func TestToInt(t *testing.T) {
	f := "toInt"
	{
		t1 := 101
		ret := checkEnv(t, gScriptContext, f, func(tf script.Function) interface{} {
			return tf.Invoke(nil, fmt.Sprintf("%d", t1))
		})
		if ret != nil && int(ret.(script.Int)) != t1 {
			t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), t1)
			return
		}
	}
	{
		t1 := true
		t1Int := 1
		ret := checkEnv(t, gScriptContext, f, func(tf script.Function) interface{} {
			return tf.Invoke(nil, t1)
		})
		if ret != nil && int(ret.(script.Int)) != t1Int {
			t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), t1Int)
			return
		}
	}
	{
		t1 := "1000001"
		t1Int := 1000001
		ret := checkEnv(t, gScriptContext, f, func(tf script.Function) interface{} {
			return tf.Invoke(nil, t1)
		})
		if ret != nil && int(ret.(script.Int)) != t1Int {
			t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), t1Int)
			return
		}
	}
}
