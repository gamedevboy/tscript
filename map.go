package script

type Map interface {
    SetValue(key, value Value)
    GetValue(key Value) Value
    ContainsKey(value Value) Bool
    Foreach(func(key, value interface{}) bool) Int
    Set(key, value interface{})
    Get(key interface{}) interface{}
    Len() Int
    Delete(key Value)
}
