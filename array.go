package script

type Array interface {
    GetSlice() []Value
    GetElement(index Int) Value
    SetElement(index Int, value Value)
    Push(args ...Value) interface{}
    Pop() interface{}
    Unshift(args ...Value) interface{}
    Shift() interface{}
    First() interface{}
    Last() interface{}
    Len() Int
    RemoveAt(index Int) Bool
    Clear()
}
