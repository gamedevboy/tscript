package map_test

import (
	"testing"
	"tklibs/script"
	"tklibs/script/compiler/test"
	"tklibs/script/runtime/context"
	_map "tklibs/script/type/map"
)

var scriptTest = `
function json_print(m) {
	return json.encode(m)
}
function forEach(m) {
    var count = 0
	m.forEach(function (k,v) {
		count++
	})
	return count
}

function setk1k2(m) {
	m["k1"] = 1
	m["k2"] = 2
	return 3
}

function length(m) {
	return m.length()
}

function delete(m,k) {
    v = m[k]
	m.delete(k)
	return v
}

function get1(m,k) {
	return m[k]
}
function get2(m,k) {
	return m.get(k)
}

function containsKey(m,k) {
	if(m.containsKey(k)) {
		return 1
	}
	return 0
}
`

var cc *context.Component

func init() {
	cc, _ = test.MustInitWithSource(scriptTest)
	cc.Run()
}

func checkEnv(t *testing.T, cc *context.Component, fName string, invoker func(script.Function) interface{}) interface{} {
	if tv := cc.ScriptGet(fName); tv == script.NullValue {
		t.Errorf("Failed:%s func invalid", fName)
		return nil
	} else {
		tf := tv.Get().(script.Function)
		ret := invoker(tf)
		if ret == nil || ret == script.NullValue {
			t.Errorf("Failed:%s return nil", fName)
			return nil
		}
		return ret
	}
}

func testIntEqual(t *testing.T, cc *context.Component, m *_map.Component, f string, excepted int, args ...interface{}) {
	ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
		args = append([]interface{}{m}, args...)
		return tf.Invoke(nil, nil, args...)
	})
	if ret != nil && int(ret.(script.Int)) != excepted {
		t.Errorf("forEach:%s got:%d expected:%d", f, int(ret.(script.Int)), excepted)
	}
}

func TestMap(t *testing.T) {
	tm := _map.NewScriptMap(cc, 0)
	tm.Set(script.String("k0"), script.Int(0))

	testIntEqual(t, cc, tm, "forEach", int(tm.Len()))
	testIntEqual(t, cc, tm, "length", int(tm.Len()))
	testIntEqual(t, cc, tm, "get1", 0, "k0")
	testIntEqual(t, cc, tm, "get2", 0, "k0")
	//testIntEqual(t, cc, tm.Component, "json_print", 0)
	testIntEqual(t, cc, tm, "setk1k2", 3)
	testIntEqual(t, cc, tm, "get1", 1, "k1")
	testIntEqual(t, cc, tm, "get2", 1, "k1")
	testIntEqual(t, cc, tm, "delete", 1, "k1")
	testIntEqual(t, cc, tm, "delete", 2, "k2")
	testIntEqual(t, cc, tm, "length", 1)
	testIntEqual(t, cc, tm, "containsKey", 1, "k0")
	testIntEqual(t, cc, tm, "containsKey", 0, "k1")
}
