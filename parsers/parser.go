package parsers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

// Parse test
func Parse() {
	src := `
	package main

	var a = 3

	func main() {
		b := 6
		var c int
	}
	`

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "", src, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.GenDecl, *ast.ValueSpec, *ast.AssignStmt:
			ast.Fprint(os.Stdout, fs, n, nil)
		}
		return true
	})

	ast.Fprint(os.Stdout, fs, f, nil)
}
