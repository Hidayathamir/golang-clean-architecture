// Command checkrules enforces the project rules that can be checked
// mechanically: codegen freshness (via CI), decorator pattern, and file naming.
//
// Rules enforced here:
//   - Decorator pattern (outbound): every outbound interface with a real
//     implementation must have a sibling {@name}_mw_logger.go defining a
//     {@Name}MwLogger type.
//   - Decorator pattern (usecase): every {@name}_usecase.go interface must have
//     a sibling {@name}_usecase_mw_logger.go defining a {@Name}MwLogger type.
//   - Decorator body: every error-returning method in a decorator must start a
//     span, retry the real call, record the error on the span, and emit the log.
//   - File naming: Go files use snake_case; decorators end in _mw_logger.go;
//     tests end in _test.go; per-package fixtures are helper_test.go.
//   - Test packages: test files under internal/ use an external test package
//     (package {dir}_test).
//   - Migrations: files under db/migrations/ are <14-digit-timestamp>_<name>.sql
//     and carry both "-- +migrate Up" and "-- +migrate Down" sections, Up first.
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

// validMigrationName matches <14-digit timestamp>_<snake_case name>.sql.
var validMigrationName = regexp.MustCompile(`^\d{14}_[a-z0-9]+(_[a-z0-9]+)*\.sql$`)

// decoratorRule describes where decorators live and how the decorator file and
// type name are derived from the interface file.
type decoratorRule struct {
	root       string                   // directory to walk
	qualifies  func(path string) bool   // is this the interface file?
	mwPath     func(path string) string // sibling decorator path
	mwLabel    func(path string) string // human-readable decorator name
	needsRetry bool                     // retry marker required in method bodies
}

func main() {
	var violations []string
	violations = append(violations, checkDecorators()...)
	violations = append(violations, checkFileNaming()...)
	violations = append(violations, checkTestPackages()...)
	violations = append(violations, checkMigrations()...)

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
			needsRetry: true,
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
			needsRetry: false,
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
				continue
			}
			violations = append(violations, checkDecoratorBody(mwPath, string(mw), rule.needsRetry)...)
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

// decoratorMethodRe matches a top-level method declaration that returns error,
// capturing the method name. Only methods on the decorator struct qualify.
var decoratorMethodRe = regexp.MustCompile(`(?m)^func\s+\([^)]*\)\s+(\w+)\([^)]*\)[^{]*\berror\b[^{]*\{`)

// decoratorBodyMarkers are the steps every decorator method must perform,
// checked as substrings of the method body. The retry step applies only to
// outbound decorators (needsRetry); usecases call Next directly.
var decoratorBodyMarkers = []struct {
	label     string
	marker    string
	onlyRetry bool
}{
	{"telemetry.Start", "telemetry.Start(ctx)", false},
	{"defer span.End", "defer span.End()", false},
	{"retry", "retrykit.", true},
	{"telemetry.RecordError", "telemetry.RecordError(span, err)", false},
	{"logkit.LogMw", "logkit.LogMw(ctx", false},
}

// checkDecoratorBody verifies that every error-returning method in a decorator
// file contains the required body steps (span, retry, record error, log).
func checkDecoratorBody(path, src string, needsRetry bool) []string {
	var violations []string

	locs := decoratorMethodRe.FindAllStringSubmatchIndex(src, -1)
	for i, loc := range locs {
		name := src[loc[2]:loc[3]]
		start := loc[1] // first byte of the body
		end := len(src)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		body := src[start:end]

		for _, m := range decoratorBodyMarkers {
			if m.onlyRetry && !needsRetry {
				continue
			}
			if !strings.Contains(body, m.marker) {
				violations = append(violations,
					fmt.Sprintf("decorator body: %s method %s missing %s", path, name, m.label))
			}
		}
	}
	return violations
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

// pkgDecl matches a file's `package <name>` declaration.
var pkgDecl = regexp.MustCompile(`(?m)^package\s+(\w+)`)

// checkTestPackages enforces that test files under internal/ use an external
// test package ({dir}_test), per the testing conventions in AGENTS.md.
func checkTestPackages() []string {
	var violations []string

	err := filepath.WalkDir("internal", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "mock" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m := pkgDecl.FindSubmatch(src)
		if m == nil {
			violations = append(violations, fmt.Sprintf("no package declaration: %s", path))
			return nil
		}

		want := filepath.Base(filepath.Dir(path)) + "_test"
		if got := string(m[1]); got != want {
			violations = append(violations,
				fmt.Sprintf("bad test package: %s declares package %s (want %s)", path, got, want))
		}
		return nil
	})
	if err != nil {
		violations = append(violations, fmt.Sprintf("walk internal: %v", err))
	}
	return violations
}

// checkMigrations enforces the migration file convention in db/migrations/:
// <14-digit-timestamp>_<name>.sql with "-- +migrate Up" before "-- +migrate Down".
func checkMigrations() []string {
	var violations []string

	err := filepath.WalkDir("db/migrations", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		name := filepath.Base(path)
		if !validMigrationName.MatchString(name) {
			violations = append(violations,
				fmt.Sprintf("bad migration name: %s (want <14-digit-timestamp>_<snake_case>.sql)", path))
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(src)
		up := strings.Index(s, "-- +migrate Up")
		down := strings.Index(s, "-- +migrate Down")
		if up == -1 {
			violations = append(violations, fmt.Sprintf("missing '-- +migrate Up' section: %s", path))
		}
		if down == -1 {
			violations = append(violations, fmt.Sprintf("missing '-- +migrate Down' section: %s", path))
		}
		if up != -1 && down != -1 && up > down {
			violations = append(violations,
				fmt.Sprintf("'-- +migrate Up' must come before '-- +migrate Down': %s", path))
		}
		return nil
	})
	if err != nil {
		violations = append(violations, fmt.Sprintf("walk db/migrations: %v", err))
	}
	return violations
}
