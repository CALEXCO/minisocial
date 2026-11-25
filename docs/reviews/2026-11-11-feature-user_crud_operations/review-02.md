# Review 02: `feature/user_crud_operations` → `main` (2026-11-23)

Re-review against [review-01](./review-01.md). Commit checked: `b1de571`.

**Verdict: 🔴 Request changes.** 1 of 21 fixed, 1 partially fixed.

## 🔴 Blocking
- [x] 1. `PUT /users/{id}` now points to `UpdateName` (`routes.go:16`). Fixed.
- [ ] 2. **Partial.** `LastInsertId` is gone, so the happy path works. But `store.go:69` still discards the `sql.Result`, so `PUT /users/99999` returns **200** for a user that doesn't exist. Check `RowsAffected`.
- [ ] 3. Not fixed. `user.go:9` still serializes `Password`, the password is still plaintext, and `00003_create_users.sql:10-13` still seeds credentials.
- [ ] 4. Not fixed. `handler.go:29` (`Create`) and `handler.go:86` (`UpdateName`, after `Atoi`) still have no `return`. In `UpdateName`, a bad ID now goes on to call the service with `id = 0`.
- [ ] 5. Not fixed. `00001` is still deleted and `00003` still creates the table.
- [ ] 6. Not fixed. `main.go:32` still uses `log.Println`.

## 🟠 Should fix in this PR
- [ ] 7. Not fixed. `store.go:25,36,60,69` don't use `ctx`.
- [ ] 8. Not fixed, and slightly **worse**. `handler.go:91` now maps *every* `UpdateName` error to 404, including DB outages and unique violations. A bad `{id}` still returns 500 (`handler.go:42,86`). `err.Error()` is still sent to clients.
- [ ] 9. Not fixed. `server.go` has no timeouts and no `Shutdown`.
- [ ] 10. Not fixed. There are no `_test.go` files.
- [ ] 11. Not fixed. `config.go:24` still calls `log.Fatal` and the message on line 29 is still wrong. `.env.example` is still missing.

## 🟡 Conventions and cleanup
- [ ] 12. Not fixed (`GetById`, `NewUserHandler`).
- [ ] 13. Not fixed (`int` in `GetById`, `int64` elsewhere).
- [ ] 14. Not fixed.
- [ ] 15. Not fixed (`store.go:36,60`).
- [ ] 16. Not fixed (`json:"email"` in, `json:"mail"` out).
- [ ] 17. Not fixed (`handler.go:35,55,70`).
- [ ] 18. Not fixed. `main.go:15,30` and the commented-out `Delete` remain. **New:** an empty `Handler.Delete` stub (`handler.go:98`) that compiles and does nothing.
- [ ] 19. Not fixed (`migrate.go:20`).
- [ ] 20. Not fixed (`service.go:18`).
- [ ] 21. Not fixed. `.idea/workspace.xml` is **still tracked** (`git ls-files`). `.gitignore` doesn't untrack files that are already in the index. See [git rm --cached](https://git-scm.com/docs/git-rm#Documentation/git-rm.txt---cached).

## New findings
None beyond the notes added to #4, #8 and #18 above.
