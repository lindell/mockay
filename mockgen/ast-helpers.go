package mockgen

import "go/ast"

func interfaceName(spec *ast.TypeSpec) string {
	return spec.Name.Name
}
