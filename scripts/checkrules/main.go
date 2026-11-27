// Command checkrules enforces the project rules that can be checked
// mechanically: codegen freshness (via CI), decorator pattern, and file naming.
//
// Rules enforced here:
//   - Decorator pattern: every outbound interface with a real implementation
//     must have a sibling {@name}_mw_logger.go defining a {@Name}MwLogger type.
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

// checkDecorators verifies every outbound interface has a matching MwLogger
// decorator file that declares a {Name}MwLogger type.
func checkDecorators() []string {
	var violations []string

	err := filepath.WalkDir("internal/outbound", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		// Skip test files and the decorators themselves.
		if strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_mw_logger.go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		names := interfaceNames(string(src))
		for _, name := range names {
			mwPath := strings.TrimSuffix(path, ".go") + "_mw_logger.go"
			mw, err := os.ReadFile(mwPath)
			if err != nil {
				violations = append(violations,
					fmt.Sprintf("decorator missing: %s defines interface %s but %s does not exist",
						path, name, filepath.Base(mwPath)))
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
		violations = append(violations, fmt.Sprintf("walk internal/outbound: %v", err))
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
