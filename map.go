package script

type Map interface {
    SetValue(key, value Value)
    GetValue(key Value) Value
    ContainsKey(value Value) Bool
    Foreach(func(key, value interface{}) bool) Int
    Len() Int
    Delete(key Value)
}
