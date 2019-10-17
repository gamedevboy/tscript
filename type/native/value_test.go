package native

import (
    "fmt"
    "reflect"
    "testing"

    "tklibs/script"
)

func TestValue_SetGet(t *testing.T) {
    v := NewValue(reflect.ValueOf(1))
    v.ScriptSet("test", new(script.Value).Set(10))
    if *v.ScriptGet("test") != *new(script.Value).Set(10) {
        t.Fail()
    }
}

type fooInnter struct {
    T1 int
}

type foo struct {
    Test int
    Abc  float32

    Inner *fooInnter
}

func (f *foo) Init() {

}

func TestValue_WithNative(t *testing.T) {
    f := &foo{
        Inner: new(fooInnter),
    }
    v := NewValue(reflect.ValueOf(f))
    f.Test = 100
    f.Abc = 1.2345
    f.Inner.T1 = 10
    fmt.Println(v.ScriptGet("Test").Get())
    fmt.Println(v.ScriptGet("Abc").Get())
    fmt.Println(v.ScriptGet("Inner").Get().(script.Object).ScriptGet("T1").Get())

    v.ScriptSet("Test", new(script.Value).Set(2000))

    fmt.Println(v.ScriptGet("init").Get())

    fmt.Println(f.Test)
}
