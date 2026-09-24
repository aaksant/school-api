# school-api

A school management REST API in Go. Built on SQLite and `sqlc`.

## Requirements

- Addition of student, teacher, class, and teaching assignment entry
- Modification of student, teacher, class, and teaching assignment
- Deletion of student, teacher, class, and teaching assignment entry
- Read student, teacher, class, and teaching assignment entry
- Authentication with login and logout
- Bulk modifications

## Extra Features

### Middlewares (done)

Middlewares are composed with a `Chain` helper, which applies them in the order they are listed.

- **Rate limiting**: per-IP token bucket (`golang.org/x/time/rate`) with configurable rate and burst. Inactive visitors are cleaned up in the background. Returns `429` with a `Retry-After` header.
- **CORS**: origin allowlist with credentials support.
- **Compression**: gzip for clients that send `Accept-Encoding: gzip`.
- **HTTP parameter pollution (HPP) protection**: keeps the last value of duplicated query/body params, with a configurable whitelist.
- **Response time**: sets the `X-Response-Time` header and logs each request with its status code.
- **Security headers**: HSTS, CSP, `X-Frame-Options`, `X-Content-Type-Options`, `Referrer-Policy`, and the `Cross-Origin-*` policies.

### Auth (in progress)

- JWT authentication with role claims
- Role-based access control middleware

### Planned

- Password reset mechanisms
- Deactivate user
- Executives, subject, and users domain

## Learning Points

### Database and Migrations

- SQLite ignores foreign keys by default. They have to be enabled per connection with `?_pragma=foreign_keys(1)` in the data source name
- SQLite can't change FK actions or add a `UNIQUE` constraint in place, so both need a full table rebuild
- `homeroom_teacher_id` is set to be nullable

### API Design

- Resources and URL paths are separate things. A path like `/teachers/{teacherId}/mappings` is only a way to reach a resource, and the resource itself is defined by its own table, repository, and package. Nesting the URL doesn't mean nesting the code or the data model
- Keep resource boundaries strict (one table, one repository, one package). Creating a teacher (`POST /teachers`) is separate from assigning a teacher to a subject (`POST /teachers/{teacherId}/mappings`), and neither belongs in the other's request body
- More routes can be a sign of good modeling, since it spreads the complexity across resources. It isn't a goal in itself
- Model resources by their real DB `id`. Re-deriving identity from name pairs caused several early bugs
