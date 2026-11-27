// Command checkrules enforces the project rules that can be checked
// mechanically: codegen freshness (via CI), decorator pattern, and file naming.
//
// Rules enforced here:
//   - Decorator pattern (outbound): every outbound interface with a real
//     implementation must have a sibling {@name}_mw_logger.go defining a
//     {@Name}MwLogger type.
//   - Decorator pattern (usecase): every {@name}_usecase.go interface must have
//     a sibling {@name}_usecase_mw_logger.go defining a {@Name}MwLogger type.
//   - File naming: Go files use snake_case; decorators end in _mw_logger.go;
//     tests end in _test.go; per-package fixtures are helper_test.go.
//
// Run with `make check-rules`. Exits non-zero on any violation.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var validFileName = regexp.MustCompile(`^[a-z0-9]+(_[a-z0-9]+)*(_test)?\.go$`)

// decoratorRule describes where decorators live and how the decorator file and
// type name are derived from the interface file.
type decoratorRule struct {
	root      string                   // directory to walk
	qualifies func(path string) bool   // is this the interface file?
	mwPath    func(path string) string // sibling decorator path
	mwLabel   func(path string) string // human-readable decorator name
}

func main() {
	var violations []string
	violations = append(violations, checkDecorators()...)
	violations = append(violations, checkFileNaming()...)

	if len(violations) > 0 {
		fmt.Fprintln(os.Stderr, "❌ check-rules failed:")
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "   - %s\n", v)
		}
		os.Exit(1)
	}
	fmt.Println("✅ check-rules passed")
}

// checkDecorators verifies every outbound and usecase interface has a matching
// MwLogger decorator file that declares a {Name}MwLogger type.
func checkDecorators() []string {
	rules := []decoratorRule{
		{
			root: "internal/outbound",
			qualifies: func(path string) bool {
				return !strings.HasSuffix(path, "_test.go") &&
					!strings.HasSuffix(path, "_mw_logger.go")
			},
			mwPath: func(path string) string {
				return strings.TrimSuffix(path, ".go") + "_mw_logger.go"
			},
			mwLabel: func(path string) string {
				return filepath.Base(strings.TrimSuffix(path, ".go") + "_mw_logger.go")
			},
		},
		{
			root: "internal/usecase",
			// Usecase interfaces live in <name>_usecase.go.
			qualifies: func(path string) bool {
				return strings.HasSuffix(path, "_usecase.go")
			},
			mwPath: func(path string) string {
				return strings.TrimSuffix(path, ".go") + "_mw_logger.go"
			},
			mwLabel: func(path string) string {
				return filepath.Base(strings.TrimSuffix(path, ".go") + "_mw_logger.go")
			},
		},
	}

	var violations []string
	for _, rule := range rules {
		violations = append(violations, checkDecoratorsIn(rule)...)
	}
	return violations
}

// checkDecoratorsIn walks one root and reports interfaces missing their decorator.
func checkDecoratorsIn(rule decoratorRule) []string {
	var violations []string

	err := filepath.WalkDir(rule.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !rule.qualifies(path) {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		names := interfaceNames(string(src))
		for _, name := range names {
			mwPath := rule.mwPath(path)
			mw, err := os.ReadFile(mwPath)
			if err != nil {
				violations = append(violations,
					fmt.Sprintf("decorator missing: %s defines interface %s but %s does not exist",
						path, name, rule.mwLabel(path)))
				continue
			}
			want := name + "MwLogger"
			if !regexp.MustCompile(`type\s+` + regexp.QuoteMeta(want) + `\s+struct`).MatchString(string(mw)) {
				violations = append(violations,
					fmt.Sprintf("decorator invalid: %s does not declare type %s", mwPath, want))
			}
		}
		return nil
	})
	if err != nil {
		violations = append(violations, fmt.Sprintf("walk %s: %v", rule.root, err))
	}
	return violations
}

// interfaceNames returns the names of interfaces declared in src, ignoring the
// moq-generated mock interfaces that also carry the `interface` keyword.
func interfaceNames(src string) []string {
	re := regexp.MustCompile(`(?m)^type\s+(\w+)\s+interface\s*\{`)
	var out []string
	for _, m := range re.FindAllStringSubmatch(src, -1) {
		out = append(out, m[1])
	}
	return out
}

// checkFileNaming enforces the naming conventions across internal/, pkg/, cmd/.
func checkFileNaming() []string {
	var violations []string

	for _, root := range []string{"internal", "pkg", "cmd"} {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				// Generated mocks are exempt.
				if d.Name() == "mock" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			name := filepath.Base(path)
			if !validFileName.MatchString(name) {
				violations = append(violations,
					fmt.Sprintf("bad file name: %s (want snake_case[_test].go)", path))
			}
			return nil
		})
		if err != nil {
			violations = append(violations, fmt.Sprintf("walk %s: %v", root, err))
		}
	}
	return violations
}
