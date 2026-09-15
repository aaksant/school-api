package teachers

type Teacher struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

// id is server generated
type CreateTeacherRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

type PatchTeacherRequest struct {
	FirstName   *string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName    *string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
}
