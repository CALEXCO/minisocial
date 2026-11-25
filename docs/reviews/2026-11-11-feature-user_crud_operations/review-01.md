# Review 01: `feature/user_crud_operations` → `main` (2026-11-11)

Checklist version of the original review. The full text with explanations and links is in
[`docs/pr-review-user-crud.md`](../../pr-review-user-crud.md).

**Verdict: 🔴 Request changes**

## 🔴 Blocking
- [ ] 1. `PUT /users/{id}` is wired to `GetById` instead of `UpdateName`
- [ ] 2. `UpdateName` always fails: `LastInsertId()` isn't supported by pgx, and not-found isn't detected (`RowsAffected`)
- [ ] 3. Passwords stored in plaintext, serialized to JSON, and seeded in the migration
- [ ] 4. Missing `return` after `http.Error` (`Create` and `UpdateName`)
- [ ] 5. Migration history rewritten (`00001` deleted, `00003` added)
- [ ] 6. A failed migration doesn't stop startup (`log.Println` instead of `log.Fatal`)

## 🟠 Should fix in this PR
- [ ] 7. `ctx` is accepted but never used (`QueryRow` / `Query` / `Exec` instead of the `*Context` variants)
- [ ] 8. No domain errors: everything returns 500, and `err.Error()` is sent to the client
- [ ] 9. `http.Server` has no timeouts and no graceful shutdown
- [ ] 10. No tests
- [ ] 11. `config.LoadConfig` calls `log.Fatal` when there's no `.env`; the default-address log message is wrong; `.env.example` is missing

## 🟡 Conventions and cleanup
- [ ] 12. Initialisms and stutter: `GetById` → `GetByID`, `NewUserHandler` → `NewHandler`
- [ ] 13. ID types are inconsistent (`int` / `int64`)
- [ ] 14. Return types are inconsistent (`*User` / `User`)
- [ ] 15. `SELECT *` with a positional `Scan`
- [ ] 16. Naming mismatch: `email` / `mail` / `Mail`, and `username` / `name`
- [ ] 17. Wrong `slog` arguments (`!BADKEY`)
- [ ] 18. Leftover debug code and commented-out code
- [ ] 19. Unused `database.Down`
- [ ] 20. Weak email validation, and no normalization
- [ ] 21. `.idea/` hygiene
