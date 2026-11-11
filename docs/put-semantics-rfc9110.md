# What should `PUT` return? (RFC 9110)

The current HTTP semantics spec is **RFC 9110**. It replaced RFC 7231 in 2022.
PUT is in [RFC 9110 §9.3.4](https://www.rfc-editor.org/rfc/rfc9110#section-9.3.4). Read that section yourself; these are the points to look for:

- **Updating a resource that already exists:** the server must answer **200 (OK)** or **204 (No Content)**. Which one you pick is a design choice, and the RFC doesn't decide it for you.
- **Creating a resource that didn't exist:** the server must answer **201 (Created)**.
- **The body doesn't fit what the resource allows:** the section mentions **409 (Conflict)** and **415 (Unsupported Media Type)**.
- **What PUT means:** the body replaces the whole state of the resource. Read the part that compares PUT with POST and partial updates.

## Status code definitions (to compare 200 and 204)

- [§15.3.1 200 OK](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.1): what a 200 response is expected to contain.
- [§15.3.5 204 No Content](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.5): read why this code suits a "save" operation.
- [§15.3.2 201 Created](https://www.rfc-editor.org/rfc/rfc9110#section-15.3.2) and the [Location header (§10.2.2)](https://www.rfc-editor.org/rfc/rfc9110#section-10.2.2).
- [§15.5.5 404 Not Found](https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5) and [§15.5.10 409 Conflict](https://www.rfc-editor.org/rfc/rfc9110#section-15.5.10): your "user not found" and "duplicate name" cases.

## Partial updates

- [RFC 5789 (PATCH)](https://www.rfc-editor.org/rfc/rfc5789): read §1 and §2 on how PATCH differs from PUT.

## Two questions to answer before you choose

Both are about your own design:

1. **Your IDs come from `BIGSERIAL`, so the server assigns them.** If a client sends `PUT /users/999` and that user doesn't exist, should the server create it (as RFC 9110 allows) or return 404? Read §9.3.4 with that in mind.
2. **Your endpoint only changes `name`.** Given what §9.3.4 says PUT means, is PUT the right method, or does RFC 5789 fit better?

## Interview point

"PUT vs PATCH" and "why PUT is idempotent and POST isn't" are common REST questions.
Idempotency is in [RFC 9110 §9.2.2](https://www.rfc-editor.org/rfc/rfc9110#section-9.2.2).
