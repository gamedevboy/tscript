package script

type Map interface {
    ContainsKey(interface{}) Bool
    Foreach(func(key, value interface{}) bool) Int
    Set(interface{},interface{})
    Get(interface{}) interface{}
    Len() Int
    Delete(interface{})
}
