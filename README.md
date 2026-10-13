# school-api

A school management REST API in Go. Built on SQLite and `sqlc`.

## Requirements

- Addition of student, teacher, class, and teaching assignment entry
- Modification of student, teacher, class, and teaching assignment
- Deletion of student, teacher, class, and teaching assignment entry
- Read student, teacher, class, and teaching assignment entry
- Authentication with login, token refresh, and logout
- Bulk modifications

## Configuration

| Variable | Purpose |
|---|---|
| `JWT_SECRET` | HS256 signing key, at least 32 bytes. The server refuses to start without it |

Variables are read from the environment, and a `.env` file is loaded when present. Keep `.env` out of version control.

## Extra Features

### Middlewares (done)

Middlewares are composed with a `Chain` helper, which applies them in the order they are listed.

- **Rate limiting**: per-IP token bucket (`golang.org/x/time/rate`) with configurable rate and burst. Inactive visitors are cleaned up in the background. Returns `429` with a `Retry-After` header.
- **CORS**: origin allowlist with credentials support.
- **Compression**: gzip for clients that send `Accept-Encoding: gzip`.
- **HTTP parameter pollution (HPP) protection**: keeps the last value of duplicated query/body params, with a configurable whitelist.
- **Response time**: sets the `X-Response-Time` header and logs each request with its status code.
- **Security headers**: HSTS, CSP, `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`, and the `Cross-Origin-*` policies.
- **Request body maximum size**: wraps the whole mux with `MaxBytesReader` so a client can't exhaust memory with a huge payload. Applied once, outermost, so every current and future route is covered.

### Auth (done)

**Tokens**

- **Access token**: a short-lived JWT (HS256, 5 minutes) carrying `role` and `person_id` claims. It is verified by signature and expiry only and is never stored.
- **Refresh token**: a random 32-byte opaque string (7 days). Only its SHA-256 hash is stored, so a leaked database can't be used to mint sessions. The raw value is returned to the client once.
- **Rotation**: every `/users/refresh` revokes the presented token and issues a new one in a single transaction. Revoking is an atomic claim (`UPDATE ... WHERE revoked_at IS NULL`), so two concurrent requests with the same token can't both succeed.
- **Signout** revokes every refresh token the caller owns ("log out everywhere").

**Passwords and login**

- Passwords are hashed with bcrypt (72-byte limit enforced at validation).
- Every failed authentication returns the same `401` with the same body, so clients can't tell an unknown email from a wrong password, an inactive account, or a revoked token.
- Login runs a bcrypt compare against a dummy hash when the email doesn't exist, so response time doesn't reveal which emails have accounts.

**Accounts**

- A login (`users` row) is attached to an existing person record (student, teacher, staff, or executive). Exactly one of the four FKs is set, enforced by a `CHECK` that ties it to `role`.
- `POST /users/register` is executive-only: role and person ID must come from someone the system trusts, never from an anonymous client.

**Endpoints**

| Route | Access |
|---|---|
| `POST /users/login` | public |
| `POST /users/refresh` | public (requires a valid refresh token) |
| `POST /users/signout` | any authenticated user |
| `POST /users/register` | executive |

**Access policy (role-level)**

| Resource | Read | Write |
|---|---|---|
| classes | any authenticated user | executive |
| students | executive, teacher | executive |
| teachers | executive | executive |
| teaching assignments | executive, teacher | executive |

Protection is applied per route with `tm.Protect(handler, roles...)`, which chains `Authenticate` and then `RequireRole`. With no roles it means "any authenticated user". Routes are opt-in, so a regression test calls every route without a token and expects `401` except for the public allowlist.

**Known limitations**

- An access token stays valid until it expires (up to 5 minutes) after signout, because verification is stateless. Signout does guarantee that no new access token can be obtained.
- Access control is role-level only. Per-record rules (a teacher seeing only their own records) are not enforced everywhere.
- The first executive must be created outside the API (seed SQL), since `register` itself needs an executive token. 

### Planned

- Password reset (admin reset or self-service via emailed token) and change password
- Deactivate user (and revoke their sessions)
- Refresh-token reuse detection (revoke all sessions when a revoked token is replayed)
- Per-record access rules
- Repository tests: rotation rollback and concurrent refresh
- Executives and subject domains
- Full CRUD for the users domain

## Learning Points

### Database and Migrations

- SQLite ignores foreign keys by default. They have to be enabled per connection with `?_pragma=foreign_keys(1)` in the data source name
- SQLite can't change FK actions or add a `UNIQUE` constraint in place, so both need a full table rebuild
- `homeroom_teacher_id` is set to be nullable
- SQLite stores timestamps as `TEXT`, and defaults like `datetime('now')` use `YYYY-MM-DD HH:MM:SS`, not RFC3339. Use one layout for every time column, always in UTC, or parsing breaks and SQL string comparisons silently give wrong answers
- A multi-row `INSERT` is atomic: one unique-constraint violation rejects all the rows

### API Design

- Resources and URL paths are separate things. A path like `/teachers/{teacherId}/mappings` is only a way to reach a resource, and the resource itself is defined by its own table, repository, and package. Nesting the URL doesn't mean nesting the code or the data model
- Keep resource boundaries strict (one table, one repository, one package). Creating a
