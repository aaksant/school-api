-- name: GetClass :one
SELECT * FROM classes
WHERE id = ?
LIMIT 1;

-- name: ListClasses :many
SELECT * FROM classes
ORDER BY id;

-- name: CreateClass :one
INSERT INTO classes (name, homeroom_teacher_id)
VALUES (?, ?)
RETURNING *;

-- name: UpdateClass :one
UPDATE classes
SET name = ?, homeroom_teacher_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteClass :one
DELETE FROM classes
WHERE id = ?
RETURNING *;
