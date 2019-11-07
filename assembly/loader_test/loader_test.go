package loader_test

import (
	"testing"
	"tklibs/script"
	"tklibs/script/runtime/context"
	"tklibs/script/testing2"
)

var scriptTest = `
function loaderTest() {
    logger.debug("dispatcherSelect")
	return "loaderTest"
}
`
var scriptTestReload = `
function loaderTest() {
    logger.debug("dispatcherSelect")
    logger.debug("dispatcherSelect reload")
	return "loaderTest"
}
`

var cc *context.Component

func init() {
	cc,_ = testing2.MustInitWithSource(scriptTest)
	cc.Run()
}

func checkEnv(t *testing.T, cc *context.Component, fName string, invoker func(script.Function) interface{}) interface{} {
	if tv := cc.ScriptGet(fName); tv == script.NullValue {
		t.Errorf("Failed:%s func invalid", fName)
		return nil
	} else {
		tf := tv.GetFunction()
		ret := invoker(tf)
		if ret == nil {
			t.Errorf("Failed:%s return nil", fName)
			return nil
		}
		return ret
	}
}

func testWithContext(t *testing.T, cc *context.Component) {
	f := "loaderTest"
	ret := checkEnv(t, cc, f, func(tf script.Function) interface{} {
		return tf.Invoke(nil)
	})
	if ret != nil && string(ret.(script.String)) != f {
		t.Errorf("Failed:%s got:%s excepted:%s", f, ret, f)
	}
}

func TestLoadSource(t *testing.T) {
	testWithContext(t, cc)
}

func TestReloadLoadSource(t *testing.T) {
	_,acReload := testing2.MustInitWithSource(scriptTestReload)
	cc.Reload(acReload)
	testWithContext(t, cc)
}

func TestLoadBinary(t *testing.T) {
	tsb := testing2.MustCompileToTemp(scriptTest)
	cc2,_ := testing2.MustInitWithFile(tsb)
	cc2.Run()
	testWithContext(t, cc2)
}
