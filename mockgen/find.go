package mockgen

import (
	"go/ast"
)

func (f *file) findAtPosition(findFunc func(ast.Node) bool, line, column int) ast.Node {
	var found ast.Node
	abortableInspect(f.astFile, func(n ast.Node) (bool, bool) {
		if n == nil {
			return true, false
		}

		pos := f.fset.Position(n.Pos())
		end := f.fset.Position(n.End())

		if line < pos.Line || (line == pos.Line && column < pos.Column) {
			return false, true
		}

		if line < end.Line || (line == end.Line && column <= end.Column) {
			if findFunc(n) {
				found = n
				return false, true
			}
		}
		return true, false
	})
	return found
}

func abortableInspect(n ast.Node, f func(ast.Node) (bool, bool)) {
	aborted := false
	ast.Inspect(n, func(n ast.Node) bool {
		if aborted {
			return false
		}
		goDown, abort := f(n)
		if abort {
			aborted = true
			return false
		}
		return goDown
	})
}
