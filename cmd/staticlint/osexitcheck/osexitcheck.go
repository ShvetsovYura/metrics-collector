package osexitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "проверяет выхов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		pkgName := pass.Pkg.Name()
		if pkgName != "main" {
			return nil, nil
		}

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					// если это не main - не идем глубже
					return false
				}
			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						pass.Reportf(s.Sel.NamePos, "Вызов os.Exit в функции main пакета main")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
