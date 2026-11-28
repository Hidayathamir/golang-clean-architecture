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
//   - Inbound spans: every HTTP controller method starts a controller span
//     (telemetry.StartController); every Kafka consumer method starts a consumer
//     span (telemetry.StartConsumer or telemetry.StartConsumerBatch).
//   - Consumer idempotency: every ConsumeEvent* call in
//     internal/inbound/messaging/route/route_messaging.go must pass a handler
//     produced by messaging.IdempotencyHandlerSingle/Batch, unless the consumer
//     group is allowlisted as inherently idempotent.
//   - Repository table/column literals: repository implementations must not use
//     string literals for table or column names; use pkg/constant/table and
//     pkg/constant/column instead. Decorators (_mw_logger.go) are exempt because
//     their literals are log field keys, and tests (_test.go) are exempt.
//   - Error wrapping: error-returning methods in the impl files of every layer
//     (outbound/*, usecase, inbound) must wrap a handled error with
//     errkit.AddFuncName before returning it. Bare propagation of an err
//     identifier (e.g. the outer return of a DB.Transaction closure, whose inner
//     body already wrapped) is exempt. Decorators and tests are exempt.
//
// Run with `make check-rules`. Exits non-zero on any violation.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
	violations = append(violations, checkInboundSpans()...)
	violations = append(violations, checkConsumerIdempotency()...)
	violations = append(violations, checkRepositoryTableColumnLiterals()...)
	violations = append(violations, checkErrorWrapping()...)

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

// inboundSpanRule describes an inbound package whose methods must start a span.
type inboundSpanRule struct {
	root       string   // directory to walk
	markers    []string // acceptable span-start calls (any one satisfies the rule)
	markerDesc string   // human-readable form of markers
	methodRe   *regexp.Regexp
}

// checkInboundSpans verifies that every controller/consumer method starts its
// own telemetry span, matching the inbound-span convention enforced in AGENTS.md.
func checkInboundSpans() []string {
	rules := []inboundSpanRule{
		{
			root:       "internal/inbound/http",
			markers:    []string{"telemetry.StartController("},
			markerDesc: "telemetry.StartController",
			methodRe:   regexp.MustCompile(`(?m)^func\s+\([^)]*\*(\w*Controller)\)\s+(\w+)\(`),
		},
		{
			root:       "internal/inbound/messaging",
			markers:    []string{"telemetry.StartConsumer(", "telemetry.StartConsumerBatch("},
			markerDesc: "telemetry.StartConsumer or telemetry.StartConsumerBatch",
			methodRe:   regexp.MustCompile(`(?m)^func\s+\([^)]*\*(\w*Consumer)\)\s+(\w+)\(`),
		},
	}

	var violations []string
	for _, rule := range rules {
		violations = append(violations, checkInboundSpansIn(rule)...)
	}
	return violations
}

// checkInboundSpansIn walks one inbound root and reports methods missing a span.
func checkInboundSpansIn(rule inboundSpanRule) []string {
	var violations []string

	err := filepath.WalkDir(rule.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(src)

		locs := rule.methodRe.FindAllStringSubmatchIndex(s, -1)
		for i, loc := range locs {
			name := s[loc[4]:loc[5]]
			start := loc[1]
			end := len(s)
			if i+1 < len(locs) {
				end = locs[i+1][0]
			}
			body := s[start:end]

			found := false
			for _, m := range rule.markers {
				if strings.Contains(body, m) {
					found = true
					break
				}
			}
			if !found {
				violations = append(violations,
					fmt.Sprintf("inbound span: %s method %s missing %s", path, name, rule.markerDesc))
			}
		}
		return nil
	})
	if err != nil {
		violations = append(violations, fmt.Sprintf("walk %s: %v", rule.root, err))
	}
	return violations
}

// routeMessagingPath is the single file that wires Kafka consumers to handlers.
const routeMessagingPath = "internal/inbound/messaging/route/route_messaging.go"

// routeFile holds the parsed route_messaging.go for the duration of the check,
// so helpers can locate the block enclosing a given call.
var routeFile *ast.File

// idempotencyAllowlist names consumer groups whose handlers are inherently
// idempotent and therefore intentionally skip the de-dup wrapper (e.g. an
// Elasticsearch upsert that converges under repeated writes). Any new unwrapped
// consumer fails the check until it is either wrapped or added here.
var idempotencyAllowlist = map[string]bool{
	"ImageUploadedSyncSearch":      true,
	"ImageUploadedSyncSearchRetry": true,
}

// consumeEventFns are the loop functions whose handler argument must be wrapped.
var consumeEventFns = map[string]bool{
	"ConsumeEventSingle": true,
	"ConsumeEventBatch":  true,
	"ConsumeEventRetry":  true,
}

// checkConsumerIdempotency verifies every ConsumeEvent* call in
// route_messaging.go receives a handler produced by
// messaging.IdempotencyHandlerSingle/Batch, unless its consumer group is
// allowlisted as inherently idempotent.
//
// Each consumer is a wg.Go(func(){...}) block whose final statement is the
// ConsumeEvent* call; the handler is the call's last argument. A wrapped handler
// is always a call to an IdempotencyHandler* function, so this is a direct AST
// inspection rather than a fragile text match.
func checkConsumerIdempotency() []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, routeMessagingPath, nil, 0)
	if err != nil {
		return []string{fmt.Sprintf("parse %s: %v", routeMessagingPath, err)}
	}
	routeFile = file

	var violations []string

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		fnName := calleeName(call.Fun)
		if !consumeEventFns[fnName] {
			return true
		}

		if group := consumerGroupForCall(fset, file, call); idempotencyAllowlist[group] {
			return true
		}

		if isIdempotencyWrapped(call, call.Args[len(call.Args)-1]) {
			return true
		}

		pos := fset.Position(call.Pos())
		violations = append(violations,
			fmt.Sprintf("consumer idempotency: %s:%d %s passes an unwrapped handler; wrap it with messaging.IdempotencyHandlerSingle/Batch (or allowlist the consumer group)",
				routeMessagingPath, pos.Line, fnName))
		return true
	})

	return violations
}

// isIdempotencyWrapped reports whether the handler expression is a call to one
// of the IdempotencyHandler* wrapper functions. The handler is usually passed as
// a local identifier (`handler`), so an Ident is resolved back to the right-hand
// side of the assignment inside the same consumer block.
func isIdempotencyWrapped(call *ast.CallExpr, expr ast.Expr) bool {
	if call, ok := expr.(*ast.CallExpr); ok {
		return isIdempotencyWrapperName(calleeName(call.Fun))
	}
	id, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	// Find the block that contains this call, then look for `handler := <wrapped>`
	// within it. Blocks are the func literals passed to wg.Go per consumer.
	block := enclosingBlock(call)
	if block == nil {
		return false
	}
	wrapped := false
	ast.Inspect(block, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}
		lhs, ok := assign.Lhs[0].(*ast.Ident)
		if !ok || lhs.Name != id.Name {
			return true
		}
		if c, ok := assign.Rhs[0].(*ast.CallExpr); ok && isIdempotencyWrapperName(calleeName(c.Fun)) {
			wrapped = true
			return false
		}
		return true
	})
	return wrapped
}

// enclosingBlock returns the tightest *ast.BlockStmt containing call, which for
// the route file is the body of the per-consumer func literal.
func enclosingBlock(call *ast.CallExpr) *ast.BlockStmt {
	var best *ast.BlockStmt
	start := call.Pos()
	end := call.End()
	ast.Inspect(routeFile, func(n ast.Node) bool {
		blk, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		if blk.Pos() <= start && blk.End() >= end {
			// Tightest span = best candidate.
			if best == nil || (blk.End()-blk.Pos()) < (best.End()-best.Pos()) {
				best = blk
			}
		}
		return true
	})
	return best
}

// isIdempotencyWrapperName reports whether a function name is a de-dup wrapper.
func isIdempotencyWrapperName(name string) bool {
	return name == "IdempotencyHandlerSingle" || name == "IdempotencyHandlerBatch"
}

// calleeName returns the final identifier of a call's function expression,
// covering both pkg.Fn (SelectorExpr) and bare Fn (Ident).
func calleeName(fun ast.Expr) string {
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	}
	return ""
}

// consumerGroupForCall returns the consumer-group identifier for a ConsumeEvent*
// call by finding the `consumerGroup := consumergroup.X` assignment that precedes
// the call in the same function literal. The assignment always sits on an earlier
// source line than the call, so a file-wide walk keyed on line order is enough.
func consumerGroupForCall(fset *token.FileSet, file ast.Node, call *ast.CallExpr) string {
	callLine := fset.Position(call.Pos()).Line
	group := ""

	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			return true
		}
		id, ok := assign.Lhs[0].(*ast.Ident)
		if !ok || id.Name != "consumerGroup" {
			return true
		}
		line := fset.Position(assign.Pos()).Line
		if line >= callLine {
			return true
		}
		if sel, ok := assign.Rhs[0].(*ast.SelectorExpr); ok {
			group = sel.Sel.Name
		}
		return true
	})
	return group
}

// repositoryImplDir holds the repository layer's real implementations.
const repositoryImplDir = "internal/outbound/repository"

// tableConstantFile and columnConstantFile are the single source of truth for
// table and column names. The rule reads their values from source so the list
// never drifts from the constants themselves.
const (
	tableConstantFile  = "pkg/constant/table/table.go"
	columnConstantFile = "pkg/constant/column/column.go"
)

// checkRepositoryTableColumnLiterals fails when a repository implementation uses
// a string literal that equals a known table or column name, since those must
// come from pkg/constant/table and pkg/constant/column. Decorators and tests are
// skipped: decorator literals are log field keys, not SQL.
func checkRepositoryTableColumnLiterals() []string {
	names := constantStringValues(tableConstantFile)
	names = append(names, constantStringValues(columnConstantFile)...)
	// "context"/"time" style package names are not table/column names, and the
	// constants never collide with import paths, so a plain set lookup is safe.
	banned := make(map[string]bool, len(names))
	for _, n := range names {
		banned[n] = true
	}

	var violations []string

	err := filepath.WalkDir(repositoryImplDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_mw_logger.go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		file, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			violations = append(violations, fmt.Sprintf("parse %s: %v", path, perr))
			return nil
		}

		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			value, uerr := strconv.Unquote(lit.Value)
			if uerr != nil || !banned[value] {
				return true
			}
			violations = append(violations,
				fmt.Sprintf("%s:%d hardcoded table/column literal %q: use pkg/constant/table or pkg/constant/column",
					path, fset.Position(lit.Pos()).Line, value))
			return true
		})
		return nil
	})
	if err != nil {
		violations = append(violations, fmt.Sprintf("walk %s: %v", repositoryImplDir, err))
	}
	return violations
}

// constantStringValues returns the string values assigned to constants in a
// Go source file, which for the table/column packages are the name literals.
func constantStringValues(path string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	var out []string
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}
		for _, v := range spec.Values {
			lit, ok := v.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			if value, err := strconv.Unquote(lit.Value); err == nil {
				out = append(out, value)
			}
		}
		return true
	})
	return out
}

// errorWrappingRoots are the layers whose impl files must wrap handled errors.
var errorWrappingRoots = []string{
	"internal/outbound/repository",
	"internal/outbound/cache",
	"internal/outbound/messaging",
	"internal/outbound/search",
	"internal/outbound/storage",
	"internal/usecase",
	"internal/inbound",
}

// checkErrorWrapping fails when a handled error is returned without
// errkit.AddFuncName. Decorators and tests are skipped; bare propagation of an
// err identifier is allowed because the originating body already wrapped it.
func checkErrorWrapping() []string {
	var violations []string
	for _, root := range errorWrappingRoots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}
			if strings.HasSuffix(path, "_mw_logger.go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			fset := token.NewFileSet()
			file, perr := parser.ParseFile(fset, path, nil, 0)
			if perr != nil {
				violations = append(violations, fmt.Sprintf("parse %s: %v", path, perr))
				return nil
			}
			violations = append(violations, unwrappedErrorReturns(fset, path, file)...)
			return nil
		})
		if err != nil {
			violations = append(violations, fmt.Sprintf("walk %s: %v", root, err))
		}
	}
	return violations
}

// unwrappedErrorReturns reports every return inside an `if err != nil` guard
// that returns an unwrapped error value. A return of a bare err identifier is
// treated as propagation and allowed.
func unwrappedErrorReturns(fset *token.FileSet, path string, file *ast.File) []string {
	var violations []string

	ast.Inspect(file, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok || !guardsOnErrNil(ifStmt.Cond) {
			return true
		}

		ast.Inspect(ifStmt.Body, func(in ast.Node) bool {
			ret, ok := in.(*ast.ReturnStmt)
			if !ok || len(ret.Results) == 0 {
				return true
			}
			if returnIsWrapped(ret) || returnIsBarePropagation(ret) || returnDoesNotReferenceErr(ret) {
				return true
			}
			violations = append(violations,
				fmt.Sprintf("%s:%d error return without errkit.AddFuncName",
					path, fset.Position(ret.Pos()).Line))
			return true
		})
		return true
	})
	return violations
}

// guardsOnErrNil reports whether a condition is `<x> err != nil` or
// `err := <call>; err != nil`.
func guardsOnErrNil(cond ast.Expr) bool {
	bin, ok := cond.(*ast.BinaryExpr)
	if !ok || bin.Op != token.NEQ {
		return false
	}
	left, ok := bin.X.(*ast.Ident)
	if !ok || left.Name != "err" {
		return false
	}
	right, ok := bin.Y.(*ast.Ident)
	return ok && right.Name == "nil"
}

// returnIsWrapped reports whether any returned expression calls AddFuncName.
func returnIsWrapped(ret *ast.ReturnStmt) bool {
	wrapped := false
	for _, res := range ret.Results {
		ast.Inspect(res, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok && calleeName(call.Fun) == "AddFuncName" {
				wrapped = true
			}
			return true
		})
	}
	return wrapped
}

// returnIsBarePropagation reports whether every result is either the err
// identifier or a zero-value literal (nil, {}, ""), i.e. the propagation of an
// already-wrapped error alongside zero values.
func returnIsBarePropagation(ret *ast.ReturnStmt) bool {
	sawErr := false
	for _, res := range ret.Results {
		if isZeroValueExpr(res) {
			continue
		}
		id, ok := res.(*ast.Ident)
		if !ok || id.Name != "err" {
			return false
		}
		sawErr = true
	}
	return sawErr
}

// returnDoesNotReferenceErr reports whether no result mentions the err
// identifier, meaning the return is not propagating a handled error at all
// (e.g. a fail-open wrapper returning handler(...)).
func returnDoesNotReferenceErr(ret *ast.ReturnStmt) bool {
	sawErr := false
	for _, res := range ret.Results {
		ast.Inspect(res, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok && id.Name == "err" {
				sawErr = true
			}
			return true
		})
	}
	return !sawErr
}

// isZeroValueExpr reports whether an expression is a nil literal or a composite
// literal with no elements (e.g. dto.ImageResponse{}), used as a placeholder
// alongside a propagated error.
func isZeroValueExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "nil"
	case *ast.CompositeLit:
		return len(e.Elts) == 0
	}
	return false
}
