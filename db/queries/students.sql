-- name: GetStudent :one
SELECT * FROM students
WHERE students.id = ?
LIMIT 1;
-- SELECT
--     s.id,
--     s.first_name,
--     s.last_name,
--     s.email,
--     s.date_of_birth,
--     s.created_at,
--     s.updated_at,
--     c.name AS class_name
-- FROM students s
-- LEFT JOIN class_student_mapping csmap ON csmap.student_id = s.id
-- LEFT JOIN classes c ON c.id = csmap.class_id
-- WHERE s.id = ?
-- LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students;
-- SELECT
--     s.id,
--     s.first_name,
--     s.last_name,
--     s.email,
--     s.date_of_birth,
--     s.created_at,
--     s.updated_at,
--     c.name AS class_name
-- FROM students s
-- LEFT JOIN class_student_mapping csmap ON csmap.student_id = s.id
-- LEFT JOIN classes c ON c.id = csmap.class_id
-- ORDER BY s.id;

-- name: CreateStudent :one
INSERT INTO students (first_name, last_name, email, date_of_birth)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteStudent :one
DELETE FROM students
WHERE id = ?
RETURNING *;

-- name: UpdateStudent :one
UPDATE students
SET
    first_name = ?,
    last_name = ?,
    email = ?,
    date_of_birth = ?
WHERE id = ?
RETURNING *;
