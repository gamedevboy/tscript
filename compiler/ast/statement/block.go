package statement

import (
	"container/list"
)

type Block interface {
	GetStatementList() *list.List
}

