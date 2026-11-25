# Review 04: `feature/user_crud_operations` → `main` (2026-11-23)

Re-review of the uncommitted changes in `handler.go` and `store.go`.
`go build ./...` passes. `go vet ./...` reports the 3 `slog` issues from #17.

**Verdict for #2: much better, but not ticked yet.** 2a and 2b are fixed; 2c is still open.

## 🔴 Blocking
- [ ] 2. Not-found detection on `UPDATE`
  - [x] 2a. Missing `return` after the 404. Resolved: the 0-rows check moved to the store, so that branch is gone from the handler.
  - [x] 2b. `sql.Result` leaking out of the store. Fixed: `RowsAffected` is checked in `store.go:75-82` and the interface is back to returning `error`. This is the right layer.
  - [ ] 2c. **Still every error → 404** (`handler.go:91-94`). The new block at `handler.go:96-99` is **dead code**: the block above already returned whenever `err != nil`, so the 500 branch can never run. You can't fix this with a second `if err != nil`, because both branches test the same condition. The handler has no way to tell "not found" apart from "DB down". See 2d.
  - [ ] 2d. **New:** `store.go:81` creates the error inline with `errors.New(...)`. Each call creates a new value, so no caller can compare against it with `errors.Is`. This is why 2c can't be solved yet, and it's the core of #8. [Go 1.13 errors: sentinel errors](https://go.dev/blog/go1.13-errors), [errors.Is](https://pkg.go.dev/errors#Is), [Effective Go: errors](https://go.dev/doc/effective_go#errors).

## 🟡 Conventions
- [ ] N2. **New:** the error string `"User does not exists"` is capitalized (and has a typo). By convention, Go error strings are lowercase with no trailing punctuation, because they get wrapped into longer messages. [Go Code Review Comments: error strings](https://go.dev/wiki/CodeReviewComments#error-strings)
- [ ] N1. Still open: the 200 response sets `Content-Type: application/json` but has no body.

## Other findings touched by this change
- [ ] 4. **Partial.** `Create` now has its `return` (`handler.go:30`). The one after `strconv.Atoi` in `UpdateName` (`handler.go:87`) is still missing, so a non-numeric `{id}` writes a 500 and then carries on into the service with `id = 0`.
- [ ] 17. `go vet` confirms it: `handler.go:36,56,71` pass `err` without a key.

## Lesson for next time
When a caller needs to react to a specific failure (not found, conflict, invalid input), **define your own domain error once at package level** and compare it with `errors.Is`. Don't create it inline with `errors.New` at the point where it's returned. This is the fix for 2c/2d and applies to #8 in general.
