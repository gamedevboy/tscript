package script

type Array interface {
    GetElement(index Int) Value
    SetElement(index Int, value Value)
    Push(args ...Value) interface{}
    Pop() interface{}
    Unshift(args ...Value) interface{}
    Shift() interface{}
    Len() Int
    RemoveAt(index Int) Bool
}
