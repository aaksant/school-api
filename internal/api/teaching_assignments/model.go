package teachingassignments

type TeachingAssignment struct {
	Id        int `json:"id"`
	ClassId   int `json:"class_id"`
	SubjectId int `json:"subject_id"`
	TeacherId int `json:"teacher_id"`
}

type CreateTeachingAssignmentRequest struct {
	ClassId   int `json:"class_id" validate:"required"`
	SubjectId int `json:"subject_id" validate:"required"`
}
