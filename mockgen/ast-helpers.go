package mockgen

import (
	"fmt"
	"go/ast"

	"github.com/lindell/mockay/astcopy"
)

func interfaceName(spec *ast.TypeSpec) string {
	return spec.Name.Name
}

func copyFieldList(list *ast.FieldList) *ast.FieldList {
	cpList := make([]*ast.Field, len(list.List))
	varI := 1
	for i, f := range list.List {
		var names []*ast.Ident
		if len(f.Names) == 0 {
			names = []*ast.Ident{
				{
					Name: fmt.Sprintf("var%d", varI),
				},
			}
			varI++
		} else {
			names = astcopy.IdentList(f.Names)
		}

		cpList[i] = &ast.Field{
			Type:  astcopy.Expr(f.Type),
			Tag:   astcopy.BasicLit(f.Tag),
			Names: names,
		}
	}

	return &ast.FieldList{
		List: cpList,
	}
}

func argsFromParams(params *ast.FieldList) []ast.Expr {
	var expr []ast.Expr
	for _, f := range params.List {
		for _, n := range f.Names {
			expr = append(expr, astcopy.Ident(n))
		}
	}
	return expr
}

func comment(str string) *ast.CommentGroup {
	return &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: str,
			},
		},
	}
}
