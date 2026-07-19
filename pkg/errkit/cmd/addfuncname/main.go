package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: addfuncname <file.go> [file.go ...]")
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		info, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skipping %s: %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			filepath.Walk(arg, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fi.IsDir() {
					return nil
				}
				if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
					if err := processFile(path); err != nil {
						fmt.Fprintf(os.Stderr, "error processing %s: %v\n", path, err)
					}
				}
				return nil
			})
		} else if strings.HasSuffix(arg, ".go") {
			if err := processFile(arg); err != nil {
				fmt.Fprintf(os.Stderr, "error processing %s: %v\n", arg, err)
			}
		}
	}
}

func processFile(filePath string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	pkgName := file.Name.Name

	type scope struct {
		start, end int
		name       string
	}
	var scopes []scope

	ast.Inspect(file, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			name := funcDeclName(pkgName, fn)
			scopes = append(scopes, scope{
				start: fset.Position(fn.Body.Pos()).Offset,
				end:   fset.Position(fn.Body.End()).Offset,
				name:  name,
			})
		case *ast.FuncLit:
			var parentName string
			for i := len(scopes) - 1; i >= 0; i-- {
				parentName = scopes[i].name
				break
			}
			if parentName == "" {
				parentName = pkgName
			}
			scopes = append(scopes, scope{
				start: fset.Position(fn.Body.Pos()).Offset,
				end:   fset.Position(fn.Body.End()).Offset,
				name:  parentName,
			})
		}
		return true
	})

	type callToUpdate struct {
		call *ast.CallExpr
		name string
	}
	var calls []callToUpdate
	changed := false

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if !isAddFuncName(call) {
			return true
		}

		callOffset := fset.Position(call.Pos()).Offset

		var enclosingName string
		for i := len(scopes) - 1; i >= 0; i-- {
			if scopes[i].start <= callOffset && callOffset <= scopes[i].end {
				enclosingName = scopes[i].name
				break
			}
		}
		if enclosingName == "" {
			return true
		}

		switch {
		case len(call.Args) == 0:
			return true
		case len(call.Args) == 1:
			calls = append(calls, callToUpdate{call, enclosingName})
			changed = true
		case len(call.Args) >= 2:
			last := call.Args[len(call.Args)-1]
			lit, ok := last.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			existing, _ := strconv.Unquote(lit.Value)
			if existing == enclosingName {
				return true
			}
			lit.Value = strconv.Quote(enclosingName)
			changed = true
		}

		return true
	})

	if !changed {
		return nil
	}

	for _, c := range calls {
		c.call.Args = append(c.call.Args, &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(c.name),
		})
	}

	var buf strings.Builder
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(buf.String()), 0644)
}

func isAddFuncName(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	return sel.Sel.Name == "AddFuncName"
}

func funcDeclName(pkgName string, fd *ast.FuncDecl) string {
	if fd.Recv != nil && len(fd.Recv.List) > 0 {
		return fmt.Sprintf("%s.(%s).%s", pkgName, recvTypeStr(fd.Recv.List[0].Type), fd.Name.Name)
	}
	return fmt.Sprintf("%s.%s", pkgName, fd.Name.Name)
}

func recvTypeStr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return "*" + recvTypeStr(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return recvTypeStr(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}
