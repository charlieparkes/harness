# Harness

## Structure

- `docs/canonical/`: accepted, approved project documentation
- `cmd/harness`
- `internal/harness`: business logic, core of application
- `internal/harness/model`: models
- `internal/harness/ui`: definition of TUI
- `internal/store`: store implementations

## Available Commands

- `make`: run make targets: generate, mocks, lint, test, and build
- `make generate`: run `go generate ./...`
- `make mocks`: generate mocks with mockery
- `make lint`: run `golangci-lint run --fix ./...`
- `make test`: run `go test ./...`
- `make build`: generate code, then compile `./bin/harness`
- `make install`: copy `./bin/harness` to `$(HOME)/.local/bin`

## Available Skills

- `create-documentation`: write new documentation
- `modify-documentation`: modify existing documentation
- `promote-documentation`: canonicalize new or modified documentation
- `my-plan`: define consistent standard for writing plans
- `my-plan-review`: define consistent standard for reviewing plans
- `my-review`: review code changes against the original plan used to make them
- `my-simple-english`: write or rewrite technical text with ASD-STE100 Simplified Technical English

## Rules

- Never run git commands without explicit user instruction.

## Feature Development

- Always run `make` (default target) after making changes to validate the build.
- `harness` follows an onion architecture. Each layer of the onion should define interfaces for the objects it needs from other layers.
- `internal/store` and `internal/harness` should not import code from local packages outside of themselves.
- When defining an interface, confirm a mock is generated via .mockery.yml
- Never assign a type to an interface as a global, blank identifier. The package which defines the interface should write a test which passes a concrete version of the implemented interface to each consumer.

## Testing

- Every top level test method should start with a test size assertion. `testsize.Small`, `testsize.Medium`, or `testsize.Large`.
- Test size assertion should always be an independent first or second line of the test (after t.Parallel). `ctx := testsize.Small(t)`, for example.
- Tests should be organized into files by size.
    - small: example_test.go
    - medium: example_medium_test.go
    - large: example_large_test.go