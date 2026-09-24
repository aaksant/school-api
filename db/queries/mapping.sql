-- name: RemoveStudentFromClass :execresult
DELETE FROM class_student_mapping
WHERE student_id = ?;


-- name: GetTeachingAssignmentById :one
SELECT * FROM teaching_assignments
WHERE id = ?;

-- name: ListTeachingAssignments :many
SELECT * FROM teaching_assignments
WHERE teacher_id = ?
ORDER BY id;

-- name: CreateTeacherAssignment :one
INSERT INTO teaching_assignments (class_id, subject_id, teacher_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteTeacherAssignment :one
DELETE FROM teaching_assignments
WHERE id  = ?
RETURNING *;

-- might be useful later
-- -- name: GetStudentMappingById :one
-- SELECT * FROM class_student_mapping
-- WHERE id = ?;

-- -- name: ListStudentMappings :many
-- SELECT * FROM class_student_mapping
-- WHERE student_id = ?
-- ORDER BY id;

-- -- name: AssignStudentToClass :exec
-- INSERT INTO class_student_mapping (student_id, class_id)
-- VALUES (?, ?)
-- -- upsert
-- ON CONFLICT(student_id) DO UPDATE SET class_id = exlcuded.class_id;
