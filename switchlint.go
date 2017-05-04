package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	flag.Parse()

	exitCode := 0
	for _, f := range flag.Args() {
		if err := inspect(f); err != nil {
			log.Print(err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func inspect(file string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		//no default case for irony and so it can be tested on itself
		switch n := n.(type) {
		case *ast.TypeSwitchStmt:
			check(fset.Position(n.Switch), n.Body)
		case *ast.SwitchStmt:
			if n.Tag != nil {
				check(fset.Position(n.Switch), n.Body)
			}
		}
		return true
	})

	return nil
}

func check(p token.Position, body *ast.BlockStmt) {
	for _, block := range body.List {
		if block.(*ast.CaseClause).List == nil {
			return
		}
	}
	fmt.Printf("%s: no default case\n", p)
}
