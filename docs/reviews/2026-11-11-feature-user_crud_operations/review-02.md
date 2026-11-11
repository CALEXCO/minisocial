# Review 02: `feature/user_crud_operations` → `main`

- **Date:** 2026-11-11
- **Reviewed at:** `5a387fa` (branch tip, includes `de6cf70` "Change PUT operation after review")
- **Verdict:** 🔴 Request changes. 1 of 21 findings from review 01 is fixed; 5 new findings.
- **Checked with:** `git diff main...feature/user_crud_operations`, `go vet ./...`

---

## Status of review 01 findings

### 🔴 Blocking

- [x] 1. **`PUT /users/{id}` is wired to the wrong handler.** Fixed in `de6cf70`: `routes.go:16` now points to `UpdateName`.
- [ ] 2. **`UpdateName` fails or lies at runtime.** *Partially fixed.* The `LastInsertId()` call is gone, so the endpoint no longer fails every time. But `store.go:69` now throws away the `sql.Result`, so `PUT /users/999` on a user that doesn't exist returns **200**. The "not found" case still isn't detected. Read what `RowsAffected` returns: [sql.Result](https://pkg.go.dev/database/sql#Result).
- [ ] 3. **Plaintext passwords, returned in every response, seeded in git.** Not changed. `user.go` still has `json:"password"`, and `00003_create_users.sql` still seeds credentials. See also #25.
- [ ] 4. **Missing `return` after `http.Error`.** Not changed. Both places are still there: `handler.go:28-30` (Create) and `handler.go:85-87` (UpdateName). In UpdateName, `PUT /users/abc` writes a 500, then calls the service with `id = 0`, then tries to write a second status. [http.Error](https://pkg.go.dev/net/http#Error)
- [ ] 5. **Rewrote migration history.** Not changed. `00001` is still deleted and `00003` still contains both the schema and the seed data.
- [ ] 6. **A failed migration doesn't stop startup.** Not changed (`main.go`, `log.Println`).

### 🟠 Should fix in this PR

- [ ] 7. **`ctx` is unused in the store.** Not changed. Every query still uses the non-`Context` variant.
- [ ] 8. **No domain errors.** Not changed, and UpdateName made it worse: see #22.
- [ ] 9. **No server timeouts and no graceful shutdown.** Not changed (`server.go`).
- [ ] 10. **No tests.** Not changed. With tests, #2, #4 and #22 would have been caught.
- [ ] 11. **`LoadConfig` calls `log.Fatal` when there's no `.env`.** Not changed. There's still no `.env.example`.

### 🟡 Conventions and cleanup

- [ ] 12. **Initialisms and stutter** (`GetById`, `NewUserHandler`). Not changed.
- [ ] 13. **Inconsistent ID types.** Not changed. `GetById` takes `int`, `UpdateName` takes `int64`, and `Create` scans into `int` and then converts.
- [ ] 14. **Inconsistent return types** (`*User` vs `User`). Not changed.
- [ ] 15. **`SELECT *` with positional `Scan`.** Not changed (`store.go:36`, `store.go:60`).
- [ ] 16. **Naming mismatch** (`email`/`mail`, `username`/`name`). Not changed.
- [ ] 17. **Wrong `slog` arguments.** Not changed. `go vet` reports it in three places: `handler.go:35`, `:55`, `:70`. [log/slog: Attrs and Values](https://pkg.go.dev/log/slog#hdr-Attrs_and_Values)
- [ ] 18. **Leftover debug and commented-out code.** Not changed. There's also a new empty `Delete` handler stub at `handler.go:98`.
- [ ] 19. **Unused `database.Down`.** Not changed.
- [ ] 20. **Weak email validation and no normalization.** Not changed.
- [ ] 21. **`.idea/`.** *Worse.* `.idea/workspace.xml` is tracked again, on this branch and on `main`: commit `30fa6f8` modified it. `.gitignore` only affects untracked files, so the entry does nothing for a file that's already in the index. Read the first paragraph of [gitignore](https://git-scm.com/docs/gitignore) and look at `--cached` in [git rm](https://git-scm.com/docs/git-rm).

---

## New findings

### 🟠 Should fix in this PR

- [ ] 22. **`UpdateName` maps every error to 404** (`handler.go:91`)
  Any error from the service becomes `404 Not Found`. That includes a lost DB connection, a cancelled context, and a duplicate-name `23505`. Combined with #2, you get the opposite of what you want: a user that really is missing returns 200, and a database outage returns 404. In K8s this is also misleading operationally, because a 404 doesn't show up as a 5xx error on your ingress or SLO dashboards. Fix it together with #8 by mapping specific domain errors with `errors.Is`, and fall back to 500.
  Read: [errors.Is](https://pkg.go.dev/errors#Is), [RFC 9110 §15.5.5 404](https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5), [§15.5.10 409](https://www.rfc-editor.org/rfc/rfc9110#section-15.5.10)

- [ ] 23. **`PUT` answers 200 with a JSON content type and an empty body** (`handler.go:94-95`)
  You removed the body after reading RFC 9110, which is a valid choice. But your response still declares `Content-Type: application/json` and sends 200 with nothing in it, and an empty body isn't valid JSON. Your own notes in `docs/put-semantics-rfc9110.md` point to the code that means "success, no content". Either send the updated representation with 200, or send no body (and no JSON content type) with the other code.
  Read: [RFC 9110 §15.3.1 200](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.1), [§15.3.5 204](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.5)

- [ ] 24. **No validation of `username`** (`service.go`, `Create` and `UpdateName`)
  An empty or whitespace-only name passes the service and reaches the DB. `NOT NULL` accepts `''`. CLAUDE.md puts validation and normalization in the service layer, but right now the only rule is the `@` check on email. Decide what a valid name is (length, allowed characters, trimming) and return a domain error such as `ErrInvalidName` when it fails.
  Read: [strings.TrimSpace](https://pkg.go.dev/strings#TrimSpace), [utf8.RuneCountInString](https://pkg.go.dev/unicode/utf8#RuneCountInString) (`len` counts bytes, not characters), [Postgres CHECK constraints](https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-CHECK-CONSTRAINTS)

- [ ] 25. **`POST /users` creates users with an empty password** (`service.go` Create → `store.go` Create)
  The handler never reads a password, so the service builds a `User` with `Password: ""` and the store inserts it. `password TEXT NOT NULL` doesn't stop this, because `''` isn't `NULL`. So every user created through the API has an empty password. Decide whether the password is part of the create request. If it is, it gets hashed (#3). If it isn't, the column shouldn't be `NOT NULL` with fake values.
  Read: [Postgres NOT NULL constraints](https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-NOT-NULL)

### 🟡 Conventions and cleanup

- [ ] 26. **`GetAll` returns a partial slice together with an error** (`store.go:46`)
  A `Scan` error returns `users, err`, but a `rows.Err()` error returns `nil, err`. The Go convention is that a function returning a non-nil error makes no promises about the other values, and callers shouldn't use them. Return the zero value consistently so nobody is tempted to.
  Read: [Effective Go: errors](https://go.dev/doc/effective_go#errors), [database/sql: Rows.Next](https://pkg.go.dev/database/sql#Rows.Next)

---

## Note on how this PR was updated

`de6cf70` fixes #1, but the commit message ("following RFC oficial doc") only mentions the PUT semantics change. When you push fixes for a review, it helps reviewers if each commit says which finding it addresses (for example, `Fix #1: route PUT to UpdateName`). That habit carries straight over to team PR workflows.
