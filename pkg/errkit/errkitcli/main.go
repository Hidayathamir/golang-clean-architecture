package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// This script will auto generate stack string for go file that calling errkit.AddFuncName.
// For example if file internal/usecase/address/update.go func Update
// calling errkit.AddFuncName("", err). By running this script it will auto generate to
// errkit.AddFuncName("address.(*AddressUsecaseImpl).Update", err).

func main() {
	target := "." // all .go file
	ignoreList := []string{
		"api/",
		"internal/mock",
	}

	err := filepath.Walk(target, func(path string, info fs.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if info.IsDir() {
			return nil
		}

		for _, ignore := range ignoreList {
			if strings.HasPrefix(path, ignore) {
				return nil
			}
		}

		isGoFile := strings.HasSuffix(path, ".go")
		isTestFile := strings.HasSuffix(path, "_test.go")
		if !isGoFile || isTestFile {
			return nil
		}

		return process(path)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	PACKAGE_NAME = "errkit"
	FUNC_NAME    = "AddFuncName"
)

func process(path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parse file '%s': %w", path, err)
	}

	isChanged := false

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		funcName := getFuncName(f, funcDecl)

		// there is so many checking in below code block. I suggest you check with
		// ast.Print(fset, funcDecl.Body) to understand deeper.
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok || len(callExpr.Args) < 1 {
				return true
			}
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok || selectorExpr.Sel.Name != FUNC_NAME {
				return true
			}
			pkg, ok := selectorExpr.X.(*ast.Ident)
			if !ok || pkg.Name != PACKAGE_NAME {
				return true
			}
			basicList, ok := callExpr.Args[0].(*ast.BasicLit)
			if ok && basicList.Value == strconv.Quote(funcName) {
				return true
			}

			// we change argument of errkit.AddFuncName to funcName
			callExpr.Args[0] = &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(funcName)}
			isChanged = true

			return true
		})
	}

	if !isChanged {
		return nil
	}

	fmt.Println("updating", path)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return fmt.Errorf("error format node: %w", err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("error write file '%s': %w", path, err)
	}

	return nil
}

func getFuncName(f *ast.File, funcDecl *ast.FuncDecl) string {
	pkgName := f.Name.Name
	funcName := funcDecl.Name.Name // func name that calling errkit.AddFuncName

	isFuncHasReceiver := funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0
	if isFuncHasReceiver {
		switch t := funcDecl.Recv.List[0].Type.(type) {
		case *ast.StarExpr: // pointer receiver
			if ident, ok := t.X.(*ast.Ident); ok {
				return fmt.Sprintf("%s.(*%s).%s", pkgName, ident.Name, funcName)
			}
		case *ast.Ident: // value receiver
			return fmt.Sprintf("%s.%s.%s", pkgName, t.Name, funcName)
		}
	}
	return fmt.Sprintf("%s.%s", pkgName, funcName)
}
