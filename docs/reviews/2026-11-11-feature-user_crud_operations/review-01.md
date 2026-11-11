# Review 01: `feature/user_crud_operations` → `main`

- **Date:** 2026-11-11
- **Reviewed at:** before `de6cf70` (the PUT route still pointed to `GetById`)
- **Verdict:** 🔴 Request changes
- **Source:** converted from `docs/pr-review-user-crud.md` into the checklist format that CLAUDE.md requires. The content is the same.

The layering is right: the `Store` interface, a service that depends on the interface, and a handler that depends on the service. Moving `PostgresStore` out of `handler.go` and dropping `DefaultServeMux` were good calls. But two of the advertised endpoints don't work, and there are security problems I can't approve.

---

## 🔴 Blocking

- [ ] 1. **`PUT /users/{id}` is wired to the wrong handler**
  `routes.go:16`: `PUT` points to `users.GetById`, so `UpdateName` can't be reached. A single request against the endpoint would have caught this. That's an argument for #10.

- [ ] 2. **`UpdateName` always fails at runtime**
  `store.go`: the pgx stdlib driver doesn't implement `LastInsertId()`, so it returns an error every time. Also, an UPDATE that matches no rows isn't an error to `database/sql`. That is exactly where your "not found" case comes from.
  Read: [sql.Result](https://pkg.go.dev/database/sql#Result) (what `RowsAffected` gives you), [Postgres RETURNING](https://www.postgresql.org/docs/current/dml-returning.html)

- [ ] 3. **Passwords are stored in plaintext and returned in every response**
  - The `User` struct serializes `Password` to JSON, so `GET /users` leaks every user's password.
  - The password is stored as plain text.
  - The migration seeds real-looking credentials into git history.

  Read: [encoding/json, struct tags](https://pkg.go.dev/encoding/json#Marshal) (look for `"-"`), [x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt), [OWASP Password Storage](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
  **Interview point:** separating the domain/persistence model from the API response model (DTO).

- [ ] 4. **`Create` handler: missing `return` after the error**
  After `http.Error(... 500)` execution carries on. It then writes a 201 ("superfluous WriteHeader" in the logs) and encodes a `nil` user. `UpdateName` has the same bug after `strconv.Atoi`.
  Read: [http.Error](https://pkg.go.dev/net/http#Error) (it doesn't end the request)

- [ ] 5. **You rewrote migration history**
  You deleted `00001_create_users.sql` and added `00003`. On any DB where `00001` already ran, goose sees version 1 applied, runs 3, and fails with "relation users already exists". Schema changes go in a new migration (`ALTER TABLE`). Seed data goes in its own migration or a dev-only script.
  Read: [goose: SQL migration files](https://pressly.github.io/goose/blog/2022/overview-sql-file/), [goose: CLI commands / allow missing](https://pressly.github.io/goose/documentation/cli-commands/)

- [ ] 6. **A failed migration doesn't stop startup**
  `main.go`: `log.Fatal` became `log.Println` on `Migrate`, so the app starts against a schema it doesn't match. In K8s you want this to crash so the pod never becomes Ready.

## 🟠 Should fix in this PR

- [ ] 7. **`ctx` is accepted everywhere and used nowhere**
  Every store method takes `ctx` but calls `QueryRow`, `Query` and `Exec`, so a client disconnect or timeout won't cancel the query.
  Read: [DB.QueryRowContext](https://pkg.go.dev/database/sql#DB.QueryRowContext), [go.dev/blog/context](https://go.dev/blog/context)

- [ ] 8. **No domain errors; everything is a 500**
  - Nonexistent ID → 500 (should be 404: `sql.ErrNoRows` → `ErrNotFound`).
  - Duplicate name/email → 500 (should be 409: `23505` → `ErrConflict`).
  - Invalid email → 500 (should be 400/422), created inline with `errors.New` so it can't be compared with `errors.Is`.
  - Non-numeric `{id}` → 500 (should be 400).
  - `err.Error()` is sent to the client, which leaks SQL details.

  Read: [Go 1.13 errors](https://go.dev/blog/go1.13-errors), [errors.As](https://pkg.go.dev/errors#As), [pgconn.PgError](https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#PgError), [Postgres error codes](https://www.postgresql.org/docs/current/errcodes-appendix.html), [RFC 9110 status codes](https://www.rfc-editor.org/rfc/rfc9110#name-status-codes)

- [ ] 9. **`http.Server` has no timeouts and no graceful shutdown**
  Without `ReadHeaderTimeout` you're exposed to Slowloris. Without shutdown handling, a SIGTERM from a rolling update drops in-flight requests.
  Read: [http.Server](https://pkg.go.dev/net/http#Server), [Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown), [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext), [Cloudflare: net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)

- [ ] 10. **No tests**
  The `Store` interface already makes the service testable with a fake, and handlers can be tested without a network.
  Read: [testing](https://pkg.go.dev/testing), [net/http/httptest](https://pkg.go.dev/net/http/httptest), [table-driven tests](https://go.dev/wiki/TableDrivenTests), [testcontainers-go Postgres](https://golang.testcontainers.org/modules/postgres/)

- [ ] 11. **`config.LoadConfig` crashes with no `.env` file**
  A missing `.env` → `log.Fatal`, which kills the app in a container where config comes from real env vars. Library code shouldn't `log.Fatal`. The default-address log says `http://0.0.0.0:8080` but sets `:8080`. `.env.example` is missing.
  Read: [godotenv](https://github.com/joho/godotenv), [12-Factor: Config](https://12factor.net/config)

## 🟡 Conventions and cleanup

- [ ] 12. **Initialisms and stutter:** `GetById` → `GetByID`, `NewUserHandler` → `NewHandler`. [Effective Go: names](https://go.dev/doc/effective_go#names), [Code Review Comments: initialisms](https://go.dev/wiki/CodeReviewComments#initialisms)
- [ ] 13. **Inconsistent ID types:** `int` vs `int64`. Pick `int64` to match `BIGSERIAL`. [strconv.ParseInt](https://pkg.go.dev/strconv#ParseInt)
- [ ] 14. **Inconsistent return types:** `Create` returns `*User`, the rest return `User`. [Code Review Comments: receiver type](https://go.dev/wiki/CodeReviewComments#receiver-type)
- [ ] 15. **`SELECT *` with positional `Scan`:** adding a column breaks every scan. Name the columns.
- [ ] 16. **Naming mismatch:** `email` (input) / `mail` (output) / `Mail` (struct) / `email` (column), and `username` vs `name`.
- [ ] 17. **Wrong `slog` arguments:** `slog.Default().Error("json encode failed", err)` → `!BADKEY`. Run `go vet`. [log/slog: Attrs and Values](https://pkg.go.dev/log/slog#hdr-Attrs_and_Values)
- [ ] 18. **Leftover debug code:** `fmt.Println(strings.Contains(...))`, the commented-out `// userStore :=`, and the commented-out `Delete`.
- [ ] 19. **Unused helper:** `database.Down` isn't called anywhere.
- [ ] 20. **Weak email validation and no normalization:** [net/mail.ParseAddress](https://pkg.go.dev/net/mail#ParseAddress); `Alex@...` and `alex@...` are different users.
- [ ] 21. **`.idea/`:** make sure it doesn't come back through `git add -A`.

---

## 🟢 Suggested order after the fixes

1. Fix #1–#6.
2. Domain errors and one status-mapping helper with a consistent JSON error body (#8).
3. `DELETE /users/{id}`: 204 on success, 404 if missing.
4. Request body limits: [http.MaxBytesReader](https://pkg.go.dev/net/http#MaxBytesReader), [Decoder.DisallowUnknownFields](https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields).
5. Server timeouts and graceful shutdown (#9).
6. Tests (#10).
7. Pagination on `GET /users`. A nil slice encodes as `null`, not `[]`. **Interview point:** keyset vs offset pagination. [Use the Index, Luke: no offset](https://use-the-index-luke.com/no-offset)
8. Consider `PATCH` for a name-only update. [RFC 5789](https://www.rfc-editor.org/rfc/rfc5789)

Also consider splitting the PR: the `.idea` cleanup and the CRUD feature are unrelated.
