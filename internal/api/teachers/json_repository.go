package teachers

import (
	"encoding/json"
	"os"
)

type JsonRepository struct {
	path string
}

func NewTeachersJsonRepository(path string) *JsonRepository {
	return &JsonRepository{path: path}
}

func (repo *JsonRepository) List() ([]Teacher, error) {
	data, err := os.ReadFile(repo.path)
	if err != nil {
		return nil, err
	}

	var teachers []Teacher
	if err := json.Unmarshal(data, &teachers); err != nil {
		return nil, err
	}

	return teachers, nil
}

func (repo *JsonRepository) GetById(id int) (Teacher, error) {
	Teachers, err := repo.List()
	if err != nil {
		return Teacher{}, err
	}

	for _, s := range Teachers {
		if s.Id == id {
			return s, nil
		}
	}

	return Teacher{}, ErrTeacherNotFound
}
