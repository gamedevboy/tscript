package formatter

import (
	"strings"

	"tklibs/script/compiler/ast"
)

func Format(node ast.Node) string {
	sb := &strings.Builder{}

	node.Format(0, sb)

	return sb.String()
}
