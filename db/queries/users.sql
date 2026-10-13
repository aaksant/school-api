-- name: GetById :one
SELECT * FROM users WHERE id = ?;

-- name: GetByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, role, student_id, teacher_id, staff_id, executive_id)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;
