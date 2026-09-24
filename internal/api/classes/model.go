package classes

type Class struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=3"`
	// nullable, if exists must be > 0
	HomeroomTeacherId *int `json:"homeroom_teacher_id" validate:"omitempty,gt=0"`
}

type CreateClassRequest struct {
	Name              string `json:"name" validate:"required,min=2,max=3"`
	HomeroomTeacherId *int   `json:"homeroom_teacher_id" validate:"omitempty,gt=0"`
}

type UpdateClassRequest struct {
	Name              string `json:"name" validate:"required,min=2,max=3"`
	HomeroomTeacherId *int   `json:"homeroom_teacher_id" validate:"omitempty,gt=0"`
}
