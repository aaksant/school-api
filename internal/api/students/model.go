package students

type Student struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	// dd-mm-yyyy breaks sqlite sorting
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

type CreateStudentRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

type UpdateStudentRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
}

type PatchStudentRequest struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
}
