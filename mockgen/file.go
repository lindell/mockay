package mockgen

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type file struct {
	astFile *ast.File
	fset    *token.FileSet
	src     []byte
	path    string
}

func openFile(path string) (*file, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("file or directory does not exist")
		}
		return nil, err
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &file{
		astFile: astFile,
		fset:    fset,
		path:    path,
	}, nil
}
