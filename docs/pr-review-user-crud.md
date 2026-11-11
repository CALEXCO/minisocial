# PR Review: `feature/user_crud_operations` → `main`

**Verdict: 🔴 Request changes**

The layering is right: the `Store` interface, a service that depends on the interface, and a handler that depends on the service. Moving `PostgresStore` out of `handler.go` and dropping `DefaultServeMux` were good calls. But two of the advertised endpoints don't work, and there are security problems I can't approve. Comments are below, most important first.

---

## 🔴 Blocking

### 1. `PUT /users/{id}` is wired to the wrong handler
`routes.go:16`: `PUT` points to `users.GetById`, so `UpdateName` can't be reached. A single request against the endpoint would have caught this. That's an argument for point 10.

### 2. `UpdateName` always fails at runtime
`store.go`: the pgx stdlib driver doesn't implement `LastInsertId()`, so it returns an error every time. Also, an UPDATE that matches no rows isn't an error to `database/sql`. That is exactly where your "not found" case comes from.
- [sql.Result](https://pkg.go.dev/database/sql#Result): read what `RowsAffected` gives you.
- [Postgres RETURNING](https://www.postgresql.org/docs/current/dml-returning.html)

### 3. Passwords are stored in plaintext and returned in every response
- The `User` struct serializes `Password` to JSON, so `GET /users` leaks every user's password.
- The password is stored as plain text.
- The migration seeds real-looking credentials into git history.

Read: [encoding/json, struct tags](https://pkg.go.dev/encoding/json#Marshal) (look for `"-"`), [x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt), [OWASP Password Storage](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html).

**Interview point:** separating the domain/persistence model from the API response model (DTO) is a classic design question.

### 4. `Create` handler: missing `return` after the error
`handler.go`: after `http.Error(... 500)` execution carries on. It then writes a 201 header ("superfluous WriteHeader" in the logs) and encodes a `nil` user. [http.Error docs](https://pkg.go.dev/net/http#Error) say plainly that it doesn't end the request. `UpdateName` has the same bug after `strconv.Atoi`.

### 5. You rewrote migration history
You deleted `00001_create_users.sql` and added `00003`. That breaks your own rule in CLAUDE.md. On any DB where `00001` already ran, goose sees version 1 applied, runs 3, and fails with "relation users already exists". In production, this is the kind of thing that breaks an Argo CD sync at 3 a.m. Schema changes go in a new migration (`ALTER TABLE`), and seed data goes in its own migration or a dev-only script.
- [goose: migration versioning](https://pressly.github.io/goose/blog/2022/overview-sql-file/)
- [goose: allow missing / out-of-order](https://pressly.github.io/goose/documentation/cli-commands/)

### 6. A failed migration doesn't stop startup
`main.go`: you changed `log.Fatal` to `log.Println` on `Migrate`, so the app starts against a schema it doesn't match. In K8s you want this to crash so the pod never becomes Ready.

---

## 🟠 Should fix in this PR

### 7. `ctx` is accepted everywhere and used nowhere
Every store method takes `ctx` but calls `QueryRow`, `Query` and `Exec`. A client disconnect or request timeout therefore won't cancel the query.
[DB.QueryRowContext](https://pkg.go.dev/database/sql#DB.QueryRowContext), [go.dev/blog/context](https://go.dev/blog/context)

### 8. No domain errors; everything is a 500
CLAUDE.md says the store translates DB errors into domain errors and the handler maps those to status codes. Right now:
- A nonexistent ID returns **500**. It should be 404: `sql.ErrNoRows` → `ErrNotFound`.
- A duplicate name or email (`UNIQUE`) returns **500**. It should be 409: Postgres code `23505` → `ErrConflict`.
- An invalid email returns **500**. It should be 400 or 422. You also create it inline with `errors.New`, so callers can't compare it with `errors.Is`.
- A non-numeric `{id}` returns **500**. It should be 400.
- `err.Error()` is sent to the client, which leaks internal or SQL details.

Read: [Working with errors in Go 1.13](https://go.dev/blog/go1.13-errors), [errors.As](https://pkg.go.dev/errors#As), [pgconn.PgError](https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#PgError), [Postgres error codes](https://www.postgresql.org/docs/current/errcodes-appendix.html), [RFC 9110 status codes](https://www.rfc-editor.org/rfc/rfc9110#name-status-codes).

### 9. `http.Server` has no timeouts and no graceful shutdown
Both are in your CLAUDE.md conventions. Without `ReadHeaderTimeout` you're exposed to Slowloris. Without shutdown handling, a SIGTERM from a K8s rolling update drops in-flight requests.
- [http.Server fields](https://pkg.go.dev/net/http#Server)
- [Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown)
- [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext)
- [Cloudflare: Go net/http timeouts](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/)

### 10. No tests
For a portfolio repo, this is the first thing a reviewer looks for. Your `Store` interface already makes the service testable with a fake, and handlers can be tested without a network.
- [testing](https://pkg.go.dev/testing), [net/http/httptest](https://pkg.go.dev/net/http/httptest)
- [Go wiki: table-driven tests](https://go.dev/wiki/TableDrivenTests)
- Integration tests against real Postgres: [testcontainers-go](https://golang.testcontainers.org/modules/postgres/)

### 11. `config.LoadConfig` crashes with no `.env` file
`godotenv.Load()` failing → `log.Fatal`. In a container or K8s, config comes from real env vars and there's no `.env`, so the app dies at boot. Treat a missing `.env` as optional. Also, library code shouldn't `log.Fatal`; return the error and let `main` decide. Your default-address log message says `http://0.0.0.0:8080`, but the value it sets is `:8080`. Also, `.env.example` is missing from the repo (another CLAUDE.md convention).
[godotenv README](https://github.com/joho/godotenv), [The Twelve-Factor App: Config](https://12factor.net/config)

---

## 🟡 Conventions and cleanup

12. **Initialisms:** `GetById` → `GetByID`, `NewUserHandler` → `NewHandler`, because `user.NewUserHandler` stutters. [Effective Go: names](https://go.dev/doc/effective_go#names), [Go Code Review Comments: initialisms](https://go.dev/wiki/CodeReviewComments#initialisms)
13. **Inconsistent ID types:** `int` in `GetById`, `int64` in `UpdateName` and the struct, and `int` → `int64` conversions in `Create`. Pick `int64` to match `BIGSERIAL`, and look at [strconv.ParseInt](https://pkg.go.dev/strconv#ParseInt).
14. **Inconsistent return types:** `Create` returns `*User`, while the other methods return `User`. Pick one and know why. [Go wiki: pointers vs values](https://go.dev/wiki/CodeReviewComments#receiver-type)
15. **`SELECT *` with positional `Scan`:** the next column you add breaks every scan. Name the columns explicitly.
16. **Naming mismatch:** JSON input uses `email`, output uses `mail`, the struct uses `Mail` and the column is `email`. Clients will hate it. The same goes for `username` vs `name`.
17. **Wrong `slog` arguments:** `slog.Default().Error("json encode failed", err)` is wrong. Arguments are key-value pairs, so you'll get `!BADKEY` in the logs. Run `go vet` (it has a `slog` analyzer). [log/slog](https://pkg.go.dev/log/slog#hdr-Attrs_and_Values)
18. **Leftover debug code:** remove `fmt.Println(strings.Contains(...))` and the commented-out `// userStore :=` in `main.go`, plus the commented-out `Delete` in the service and store. Either implement it or leave it out; `git` remembers.
19. **Unused helper:** `database.Down` isn't called anywhere. Remove it or wire it to a CLI flag.
20. **Email validation:** `strings.Contains(email, "@")` is weak. Look at [net/mail.ParseAddress](https://pkg.go.dev/net/mail#ParseAddress). Also normalize case and whitespace in the service, since `Alex@...` and `alex@...` are currently different users.
21. **`.idea/` is still on disk:** it's ignored now, which is fine. Make sure it doesn't come back through `git add -A`.

---

## 🟢 What to add, in order

1. **Fix #1–#6.** Those block the merge.
2. **Domain errors and status mapping** (#8): `ErrNotFound`, `ErrConflict` and `ErrInvalidEmail`, with one helper in the handler that turns them into status codes and a consistent JSON error body.
3. **`DELETE /users/{id}`:** return 204 on success and 404 if the user doesn't exist.
4. **A request body size limit.** You decode untrusted JSON with no limit and accept unknown fields. See [http.MaxBytesReader](https://pkg.go.dev/net/http#MaxBytesReader) and [Decoder.DisallowUnknownFields](https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields).
5. **Server timeouts and graceful shutdown** (#9).
6. **Tests:** service unit tests with a fake `Store`, and handler tests with `httptest`.
7. **Pagination on `GET /users`** (`limit` and a cursor). Also, an empty table currently returns `null` instead of `[]`, because a nil slice encodes as `null`. **Interview point:** cursor (keyset) vs offset pagination. [Use the Index, Luke: pagination](https://use-the-index-luke.com/no-offset)
8. **Consider `PATCH`:** `PUT` means "replace the whole resource". For a partial update of the name, `PATCH` is the more correct verb. [RFC 5789](https://www.rfc-editor.org/rfc/rfc5789)

Also consider splitting this PR. The repo cleanup (`.idea`) and the CRUD feature are unrelated changes, and smaller PRs are easier to review.
