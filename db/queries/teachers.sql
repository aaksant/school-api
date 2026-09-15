-- name: GetTeacher :one
SELECT * FROM teachers
WHERE teachers.id = ?
LIMIT 1;

-- name: ListTeachers :many
SELECT * FROM teachers;

-- name: CreateTeacher :one
INSERT INTO teachers (first_name, last_name, email, date_of_birth)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteTeacher :one
DELETE FROM teachers
WHERE id = ?
RETURNING *;

-- name: UpdateTeacher :one
UPDATE teachers
SET
    first_name = ?,
    last_name = ?,
    email = ?,
    date_of_birth = ?
WHERE id = ?
RETURNING *;


-- TODO: revisit when classes model is wip
-- -- name: ListTeachers :many
-- SELECT
--     t.id,
--     t.first_name,
--     t.last_name,
--     t.email,
--     t.date_of_birth,
--     s.name AS subject_name,
--     c.name AS class_name,
--     t.created_at,
--     t.updated_at
-- FROM teachers t
-- LEFT JOIN teaching_assignments cstmap ON cstmap.teacher_id = t.id
-- LEFT JOIN subjects s ON s.id = cstmap.subject_id
-- LEFT JOIN classes c on c.id = cstmap.class_id
-- ORDER BY t.id;

-- -- name: ListTeacherAssignments :many
-- SELECT
--     t.id,
--     t.first_name,
--     t.last_name,
--     t.email,
--     t.date_of_birth,
--     s.name AS subject_name,
--     c.name AS class_name,
--     t.created_at,
--     t.updated_at
-- FROM teachers t
-- LEFT JOIN teaching_assignments cstmap ON cstmap.teacher_id = t.id
-- LEFT JOIN subjects s ON s.id = cstmap.subject_id
-- LEFT JOIN classes c on c.id = cstmap.class_id
-- WHERE t.id = ?
