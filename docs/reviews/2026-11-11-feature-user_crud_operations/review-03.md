# Review 03: `feature/user_crud_operations` → `main` (2026-11-23)

Re-review of #2 only, against the uncommitted changes in `handler.go`, `service.go` and `store.go`.
Every other finding keeps the status it has in [review-02](./review-02.md).

**Verdict for #2: not ticked.** The 404 decision is correct, but the implementation has a bug and breaks the layering.

## 🔴 Blocking
- [ ] 2. **Still partial.** You now read `RowsAffected` and return 404 when it's 0, which is the right design choice per RFC 9110 because the server assigns the IDs. Three problems remain:
  - [ ] 2a. **Missing `return` after the 404** (`handler.go:101-103`). Execution falls through to `WriteHeader(200)`. The client still gets 404 because the first `WriteHeader` wins, but you'll see `http: superfluous response.WriteHeader call` in the logs. This is the same bug as #4. [http.Error](https://pkg.go.dev/net/http#Error) says it does not end the request.
  - [ ] 2b. **`sql.Result` leaks out of the store** (`store.go:11`, `service.go:33`, `handler.go:95`). The handler now depends on `database/sql` semantics, and the `Store` interface is tied to `database/sql`. A fake `Store` in a unit test would have to build a `sql.Result`, and a future non-SQL store couldn't implement the interface. Per your CLAUDE.md, the *store* turns "0 rows" into a domain error such as `ErrNotFound`, and the handler only maps domain errors to status codes. [Go 1.13 errors: sentinel errors](https://go.dev/blog/go1.13-errors), [Go Code Review Comments: interfaces](https://go.dev/wiki/CodeReviewComments#interfaces).
  - [ ] 2c. **`handler.go:91` still maps every service error to 404.** A DB outage or a `UNIQUE` violation on `name` (it should be 409) is reported as "not found". That overlaps with #8, but it now directly contradicts the new 0-rows branch: two different paths both say "404" for different reasons.

## 🟡 Nit
- [ ] N1. The 200 response sets `Content-Type: application/json` but has no body. Either return the updated resource or reconsider the status. See [RFC 9110 §15.3.5 204 No Content](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.5) and your own `docs/put-semantics-rfc9110.md`.

## Interview point
"Where do you check for not-found: in the handler, the service or the repository, and why?" Leaking `sql.Result` past the repository is a classic example of a leaky abstraction.
