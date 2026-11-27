// Command checkcommitmsg enforces the Conventional Commits 1.0.0 spec
// (see .opencode/skills/git-commit-convention/SKILL.md) on a commit message.
//
// Usage: checkcommitmsg <path-to-commit-msg-file>
//
// Git passes the commit message file (COMMIT_EDITMSG) as $1 to the commit-msg
// hook. Lines beginning with '#' (comments) and empty lines are ignored.
//
// Rules enforced here:
//   - Header: <type>[(scope)][!]: <description>
//   - type is a lowercase noun from the allowed set
//   - scope, when present, is a noun in parentheses
//   - description is non-empty and separated from the colon by a space
//   - a body, when present, is separated from the header by a blank line
//   - BREAKING CHANGE, when used as a footer, is uppercase
//
// Merge, revert and fixup/squash commits are exempt (git generates them).
//
// Run with `make check-commit-msg MSG_FILE=<file>`. Exits non-zero on violation.
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
	ansiGreen = "\033[32m"
)

// allowedTypes is the conventional-commit type set. feat and fix carry SemVer
// meaning; the rest mirror @commitlint/config-conventional.
var allowedTypes = map[string]bool{
	"feat":     true,
	"fix":      true,
	"build":    true,
	"chore":    true,
	"ci":       true,
	"docs":     true,
	"style":    true,
	"refactor": true,
	"perf":     true,
	"test":     true,
	"revert":   true,
}

// header matches: type, optional (scope), optional !, then ": " and description.
var header = regexp.MustCompile(`^([a-z]+)(?:\(([^()]+)\))?(!)?: (.+)$`)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: checkcommitmsg <commit-msg-file>")
		os.Exit(2)
	}

	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", os.Args[1], err)
		os.Exit(2)
	}

	lines := stripComments(string(raw))
	if len(lines) == 0 {
		// Empty message (e.g. git aborts first) — nothing to check.
		return
	}

	subject := lines[0]

	// Git-generated messages are exempt.
	if isExempt(subject) {
		return
	}

	var violations []string
	violations = append(violations, checkHeader(subject)...)
	violations = append(violations, checkBodyAndFooters(lines)...)

	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "%s❌ commit-msg failed:%s\n", ansiRed, ansiReset)
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "   - %s\n", v)
		}
		fmt.Fprintf(os.Stderr, "\n   Expected: <type>[(scope)][!]: <description>\n")
		fmt.Fprintf(os.Stderr, "   Types: feat, fix, build, chore, ci, docs, style, refactor, perf, test, revert\n")
		fmt.Fprintf(os.Stderr, "   See .opencode/skills/git-commit-convention/SKILL.md\n")
		os.Exit(1)
	}
	fmt.Printf("%s✅ commit-msg passed%s\n", ansiGreen, ansiReset)
}

// stripComments removes comment lines and trailing blank lines, returning the
// meaningful message lines up to the first blank-line-separated footer block.
func stripComments(msg string) []string {
	var lines []string
	for _, l := range strings.Split(msg, "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.HasPrefix(strings.TrimSpace(l), "#") {
			continue
		}
		lines = append(lines, l)
	}
	// Drop trailing empty lines.
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// isExempt reports whether git generated the message rather than a human.
func isExempt(subject string) bool {
	s := strings.TrimSpace(subject)
	return strings.HasPrefix(s, "Merge ") ||
		strings.HasPrefix(s, "Revert \"") ||
		strings.HasPrefix(s, "fixup!") ||
		strings.HasPrefix(s, "squash!")
}

// checkHeader validates the first line against the header grammar.
func checkHeader(subject string) []string {
	m := header.FindStringSubmatch(subject)
	if m == nil {
		return []string{
			fmt.Sprintf("header %q does not match <type>[(scope)][!]: <description>", subject),
		}
	}

	var violations []string
	typ, scope, bang, desc := m[1], m[2], m[3], m[4]

	if !allowedTypes[typ] {
		violations = append(violations,
			fmt.Sprintf("unknown type %q (allowed: %s)", typ, typeList()))
	}
	if scope == "" && strings.Contains(subject, "()") {
		violations = append(violations, "scope is empty — remove () or name the scope")
	}
	if strings.TrimSpace(desc) == "" {
		violations = append(violations, "description is empty")
	}
	// A leading-uppercase description is allowed but lowercase is conventional.
	_ = bang
	return violations
}

// checkBodyAndFooters validates the blank-line separation and footer casing.
func checkBodyAndFooters(lines []string) []string {
	var violations []string

	// Rule 6: a body must begin one blank line after the description.
	if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
		violations = append(violations,
			"body must begin one blank line after the description")
	}

	// Rule 12/15: BREAKING CHANGE, when a footer, must be uppercase.
	for _, l := range lines[1:] {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(strings.ToLower(trimmed), "breaking change") &&
			!strings.HasPrefix(trimmed, "BREAKING CHANGE") {
			violations = append(violations,
				"footer must be uppercase: BREAKING CHANGE")
		}
	}
	return violations
}

func typeList() string {
	// Deterministic order for stable output.
	order := []string{
		"feat", "fix", "build", "chore", "ci",
		"docs", "style", "refactor", "perf", "test", "revert",
	}
	return strings.Join(order, ", ")
}
