package math_test

import (
    "fmt"
    "math"
    "testing"

    "tklibs/script"
    "tklibs/script/compiler/test"
    "tklibs/script/runtime/context"
)

var scriptTest = `
function maxInt32() {
	return  math.maxInt32()
}
function toInt(val) {
	return  math.toInt(val)
}
function max(v1,v2) {
	return  math.max(v1,v2)
}
`

var cc *context.Component

func init() {
	cc, _ = test.MustInitWithSource(scriptTest)
	cc.Run()
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
	ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
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
		ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
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
		ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
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
		ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
			return tf.Invoke(nil, t1)
		})
		if ret != nil && int(ret.(script.Int)) != t1Int {
			t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), t1Int)
			return
		}
	}
	{
		f = "max"
		t1 := 1
		t2 := 2
		ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
			return tf.Invoke(nil, t1,t2)
		})
		if ret != nil && int(ret.(script.Int)) != t2 {
			t.Errorf("Failed %s: actual: %v, expected: %v", f, int32(ret.(script.Int)), t2)
			return
		}
	}
}
