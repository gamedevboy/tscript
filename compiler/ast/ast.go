package ast

import (
	"strings"
)

type Node interface {
	Format(ident int, formatBuilder *strings.Builder)
}
