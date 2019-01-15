package mockgen

import (
	"errors"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"os"

	"github.com/go-toolsmith/astcopy"
)

// Generator does contain information what should be fixed in the code and how
type Generator struct {
	logger   Logger
	fset     token.FileSet
	position *Position
	writer   io.Writer
}

// New creates a new Generator
func New(opts ...Option) *Generator {
	f := &Generator{
		logger: &nopLogger{},
		writer: os.Stdout,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Option is an option for the mock generator
type Option func(*Generator)

// Logger is the logger used
type Logger interface {
	Info(string)
}

// nopLogger is used when no other logger is specified
type nopLogger struct{}

func (n *nopLogger) Info(string) {}

// Position is the position in a document
type Position struct {
	X int
	Y int
}

// WithLogger sets a logger of the
func WithLogger(logger Logger) Option {
	return func(f *Generator) { f.logger = logger }
}

// WithPosition sets a position
func WithPosition(pos Position) Option {
	return func(f *Generator) { f.position = &pos }
}

// WithWriter sets the priter to be used
func WithWriter(writer io.Writer) Option {
	return func(f *Generator) { f.writer = writer }
}

// Generate a mock
func (f *Generator) Generate(path string) error {
	typeSpec, err := f.findInterfaceTypeSpec(path)
	if err != nil {
		return err
	}

	var fieldList []*ast.Field
	var funcDecs []ast.Decl
	interf := typeSpec.Type.(*ast.InterfaceType)
	for _, method := range interf.Methods.List {
		name := method.Names[0].Name
		fun := method.Type.(*ast.FuncType)

		mockFunc := &ast.Field{
			Names: []*ast.Ident{
				{
					Name: name + "Func",
				},
			},
			Type: astcopy.FuncType(fun),
		}
		fieldList = append(fieldList, mockFunc)

		funcDec := &ast.FuncDecl{
			Doc: &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Text: "// " + name + " mock",
					},
				},
			},
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "m",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "Mocked",
							},
						},
					},
				},
			},
			Name: &ast.Ident{
				Name: name,
			},
			Type: astcopy.FuncType(fun),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "m",
									},
									Sel: &ast.Ident{
										Name: name + "Func",
									},
								},
								Args: []ast.Expr{},
							},
						},
					},
				},
			},
		}
		funcDecs = append(funcDecs, funcDec)
	}

	genStruct := &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: "// Mocked ...",
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "Mocked",
				},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: fieldList,
					},
				},
			},
		},
	}

	decls := []ast.Decl{genStruct}
	decls = append(decls, funcDecs...)
	file := &ast.File{
		Name: &ast.Ident{
			Name: "mock",
		},
		Decls: decls,
	}

	// fmt.Println(typeSpec)
	// ast.Print(nil, genStruct)

	fset := token.NewFileSet()
	err = format.Node(f.writer, fset, file)
	if err != nil {
		return err
	}

	return nil
}

func (f *Generator) findInterfaceTypeSpec(path string) (*ast.TypeSpec, error) {
	file, err := openFile(path)
	if err != nil {
		return nil, err
	}

	if f.position == nil {
		return nil, errors.New("did not get any position")
	}

	// ast.Print(file.fset, file.astFile)

	node := file.findAtPosition(func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return false
		}

		_, ok = spec.Type.(*ast.InterfaceType)
		return ok
	}, f.position.X, f.position.Y)
	if node == nil {
		return nil, errors.New("could not find interface")
	}
	inter := node.(*ast.TypeSpec)
	return inter, nil
}
