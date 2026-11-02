# Go code style

Use `gofmt` for syntax formatting and `.editorconfig` for text files. These rules add readable structure that `gofmt` does not choose for us.

## Function layout

- Keep one task per function. Prefer early returns over nested branches.
- Separate setup, validation, a loop or external call, and the final return with a blank line.
- Keep related assignments together; do not insert a blank line between every statement.
- Name variables after their domain meaning (`job`, `section`, `sourcePath`) rather than type or length.
- Use explicit constructors for dependencies. Avoid package globals for mutable state.

```go
func collectSections(items []Item) ([]Section, error) {
	sections := make([]Section, 0, len(items))

	for _, item := range items {
		section, err := parseSection(item)
		if err != nil {
			return nil, fmt.Errorf("parse section: %w", err)
		}

		sections = append(sections, section)
	}

	return sections, nil
}
```

## Comments, errors, and logs

- Write simple English comments only where the reason is not obvious. Usually one short line; two lines only when needed. Avoid comments that repeat the next statement.
- Wrap errors with a short action and `%w`, for example `fmt.Errorf("save job: %w", err)`. Keep user-facing messages separate from internal errors.
- Log one event at a useful boundary with job ID, stage, and outcome. Never log uploaded code, secrets, full prompts, or archive contents.
- Avoid `panic` in request and worker paths. Return errors to the caller and map them to safe user states.

## Packages and tests

- Put business rules in `internal/documents` and state rules in `internal/statemachine`; keep HTTP handlers and storage adapters thin.
- Keep interfaces near the package that consumes them. Do not make a generic `utils` package.
- Automated tests are reserved for the final milestone. When added, use table-driven tests for state transitions, validation, and pure extraction rules. Use integration tests for PostgreSQL and end-to-end tests for the full document flow.
