package statement

import "container/list"

type Switch interface {
    GetTargetValue() interface{}
    SetTargetValue(value interface{})

    GetCaseList() *list.List

    GetDefaultCase() interface{}
    SetDefaultCase(value interface{})
}
