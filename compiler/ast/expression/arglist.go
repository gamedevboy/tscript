package expression

import (
    "container/list"
)

type ArgList interface {
    GetExpressionList() *list.List
}
