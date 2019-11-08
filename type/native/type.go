package native

import (
    "fmt"
    "reflect"
    "sync"

    "tklibs/script"
    "tklibs/script/runtime/native"
)

type Type struct {
    *script.ComponentType
    initOnce   sync.Once
    nativeType reflect.Type
    fieldNames []string
    fields     []field
    methods    []function
    parent     *Type
    name       string
}

func (t *Type) AddChild(fieldName string) interface{} {
    ret := &Type{
        nativeType: t.nativeType,
        fieldNames: make([]string, len(t.fieldNames), len(t.fieldNames)+1),
        parent:     t,
        name:       t.name,
    }

    copy(ret.fieldNames, t.fieldNames)
    ret.fieldNames = append(ret.fieldNames, fieldName)

    return ret
}

func (t *Type) RemoveChild(fieldName string) interface{} {
    panic("implement me")
}

func (t *Type) GetName() string {
    return t.name
}

func (t *Type) GetParent() interface{} {
    return t.parent
}

func (t *Type) GetFieldIndexByName(fieldName string) int {
    for i, v := range t.fieldNames {
        if v == fieldName {
            return i
        }
    }

    return -1
}

func (t *Type) GetFieldNames() []string {
    return t.fieldNames
}

func (t *Type) init(nativeType reflect.Type) {
    t.initOnce.Do(func() {
        t.nativeType = nativeType.Elem()
        switch t.nativeType.Kind() {
        case reflect.Interface:
            t.initForMethods(nativeType)
        case reflect.Struct:
            t.initForFields()
            t.initForMethods(nativeType)
        default:
            panic(fmt.Errorf("Not support for %v ", nativeType))
        }
    })
}

func (t *Type) initForMethods(nt reflect.Type) {
    t.methods = make([]function, nt.NumMethod())
    for i := 0; i < nt.NumMethod(); i++ {
        method := nt.Method(i)
        t.methods[i].name = method.Name
        t.methods[i].funcValue = native.ReflectFunction(method.Func)
    }
}

func (t *Type) initForFields() {
    t.fields = make([]field, t.nativeType.NumField())
    t.fieldNames = make([]string, t.nativeType.NumField())
    for i := 0; i < t.nativeType.NumField(); i++ {
        field := t.nativeType.Field(i)
        t.fieldNames[i] = field.Name
        switch field.Type.Kind() {
        case reflect.Int, reflect.Int64:
            t.fields[i].getter = func(v *Value, index int) *script.Value {
                return new(script.Value).Set(v.value.Field(index).Int())
            }
            t.fields[i].setter = func(v *Value, index int, value *script.Value) {
                switch value.GetType() {
                case script.ValueTypeInt:
                    v.value.Field(index).SetInt(int64(value.GetInt()))
                case script.ValueTypeInterface:
                    switch val := value.GetInterface().(type) {
                    case script.Int64:
                        v.value.Field(index).SetInt(int64(val))
                    default:
                        panic(fmt.Errorf("Can not convert %v to int ", reflect.ValueOf(val).Type().Name()))
                    }
                default:
                    panic(fmt.Errorf("Can not convert %v to int ", value.GetType()))
                }
            }
        case reflect.Float32, reflect.Float64:
            t.fields[i].getter = func(v *Value, index int) *script.Value {
                return new(script.Value).Set(v.value.Field(index).Float())
            }
            t.fields[i].setter = func(v *Value, index int, value *script.Value) {
                switch value.GetType() {
                case script.ValueTypeFloat:
                    v.value.Field(index).SetFloat(float64(value.GetFloat()))
                case script.ValueTypeInterface:
                    switch val := value.GetInterface().(type) {
                    case script.Float64:
                        v.value.Field(index).SetFloat(float64(val))
                    default:
                        panic(fmt.Errorf("Can not convert %v to int ", reflect.ValueOf(val).Type().Name()))
                    }
                default:
                    panic(fmt.Errorf("Can not convert %v to int ", value.GetType()))
                }
            }
        case reflect.Bool:
            t.fields[i].getter = func(v *Value, index int) *script.Value {
                return new(script.Value).Set(v.value.Field(index).Bool())
            }
            t.fields[i].setter = func(v *Value, index int, value *script.Value) {
                switch value.GetType() {
                case script.ValueTypeBool:
                    v.value.Field(index).SetBool(bool(value.GetBool()))
                default:
                    panic(fmt.Errorf("Can not convert %v to int ", value.GetType()))
                }
            }
        case reflect.Struct:
            t.fields[i].getter = func(v *Value, index int) *script.Value {
                return new(script.Value).Set(NewValue(v.value.Field(index).Addr()))
            }
        case reflect.Ptr:
            t.fields[i].getter = func(v *Value, index int) *script.Value {
                return new(script.Value).Set(NewValue(v.value.Field(index)))
            }
        default:
            panic(fmt.Errorf("Unsupport for %v ", field.Type.Kind()))
        }
    }
}

var nativeType sync.Map

func GetType(t reflect.Type) *Type {
    ret, _ := nativeType.LoadOrStore(t, new(Type))
    _type := ret.(*Type)
    _type.init(t)
    return _type
}
