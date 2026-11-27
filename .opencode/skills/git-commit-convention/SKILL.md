---
name: git-commit-convention
description: Conventional Commits 1.0.0 — commit message format, allowed types, scopes, breaking changes, and footers. Use when writing git commit messages or reviewing commit history.
---

# Conventional Commits 1.0.0

> **Enforcement**: ✅ **Enforced**. Commit messages are validated by the commit-msg hook (`.githooks/commit-msg` → `make check-commit-msg`, `scripts/checkcommitmsg`). Header grammar, allowed types, blank-line body separation and uppercase `BREAKING CHANGE` footers are checked; merge/revert/fixup/squash messages are exempt. Bypass: `git commit --no-verify`.

A lightweight convention on top of commit messages. It gives commit history explicit, machine-readable meaning and dovetails with [SemVer](https://semver.org): `fix` → PATCH, `feat` → MINOR, any `BREAKING CHANGE` → MAJOR.

## Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Structural elements

- **`fix:`** — patches a bug (correlates with PATCH).
- **`feat:`** — introduces a new feature (correlates with MINOR).
- **`BREAKING CHANGE:`** — a footer `BREAKING CHANGE: <desc>`, or a `!` after the type/scope, signals a breaking API change (correlates with MAJOR). Can be part of any type.
- Other types are allowed. `@commitlint/config-conventional` (Angular convention) recommends: `build:`, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:`, `revert:`.
- **Scope** — optional noun in parentheses describing the section, e.g. `fix(parser):`.
- **Footers** — follow the [git trailer format](https://git-scm.com/docs/git-interpret-trailers), e.g. `Refs: #123`, `Reviewed-by: Z`.

## Rules (RFC 2119 keywords)

1. Commits **MUST** be prefixed with a type (noun), then optional scope, optional `!`, then a **required** colon and space.
2. `feat` **MUST** be used for a new feature.
3. `fix` **MUST** be used for a bug fix.
4. A scope **MAY** be given; it **MUST** be a noun in parentheses, e.g. `fix(parser):`.
5. A description **MUST** immediately follow the colon+space; it's a short summary.
6. A body **MAY** follow, **MUST** begin one blank line after the description.
7. Body is free-form: any number of newline-separated paragraphs.
8. One or more footers **MAY** follow one blank line after the body. Each footer **MUST** be a word token + (`:<space>` or `<space>#`) + value.
9. Footer tokens **MUST** use `-` instead of whitespace, e.g. `Acked-by`. Exception: `BREAKING CHANGE`.
10. Footer values **MAY** contain spaces/newlines; parsing ends at the next valid footer token.
11. Breaking changes **MUST** be indicated in the type/scope prefix **or** as a footer.
12. As a footer: `BREAKING CHANGE: <description>`, uppercase.
13. In the prefix: `!` immediately before `:`. If `!` is used, the footer **MAY** be omitted.
14. Types other than `feat`/`fix` **MAY** be used.
15. Elements are **not** case-sensitive, except `BREAKING CHANGE` which **MUST** be uppercase.
16. `BREAKING-CHANGE` is synonymous with `BREAKING CHANGE` as a footer token.

## Examples

```
feat: allow provided config object to extend other configs

BREAKING CHANGE: `extends` key in config file is now used for extending other config files
```

```
feat!: send an email to the customer when a product is shipped
```

```
feat(api)!: send an email to the customer when a product is shipped
```

```
feat!: drop support for Node 6

BREAKING CHANGE: use JavaScript features not available in Node 6.
```

```
docs: correct spelling of CHANGELOG
```

```
feat(lang): add Polish language
```

```
fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123
```

```
revert: let us never again speak of the noodle incident

Refs: 676104e, a215868
```

## Notes

- If a change fits more than one type, split it into multiple commits.
- Casing is free but be consistent; conventional tools expect the lowercase default (`feat`, not `Feat`).
- Revert behavior isn't defined by the spec; using `revert:` with a `Refs:` footer to the reverted SHAs is the common recommendation.
- Reference: <https://www.conventionalcommits.org/en/v1.0.0/>

Applies to every commit made in this repo.