package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-go-file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	fset := token.NewFileSet()

	// Parse the Go file
	parsedFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse file: %v\n", err)
		os.Exit(1)
	}

	// Filter the declarations to keep only:
	// - Type declarations
	// - Function declarations (but remove the bodies)
	parsedFile.Decls = filterDecls(parsedFile.Decls)

	// Print the modified AST to stdout
	if err := printer.Fprint(os.Stdout, fset, parsedFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error printing AST: %v\n", err)
		os.Exit(1)
	}
}

// filterDecls returns only type declarations and function signatures.
// Function bodies are set to nil, so only the signature remains.
func filterDecls(decls []ast.Decl) []ast.Decl {
	var filtered []ast.Decl

	for _, decl := range decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// Keep the GenDecl only if it's a type declaration (e.g. "type Foo struct { ... }")
			if d.Tok == token.TYPE {
				filtered = append(filtered, d)
			}

		case *ast.FuncDecl:
			// Keep the function, but remove its body
			d.Body = nil
			filtered = append(filtered, d)
		}
	}

	return filtered
}
