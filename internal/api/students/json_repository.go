package students

import (
	"encoding/json"
	"os"
	"slices"
	"sync"
)

type JsonRepository struct {
	path string
	mu   sync.Mutex
}

func NewJsonRepository(path string) *JsonRepository {
	return &JsonRepository{
		path: path,
		mu:   sync.Mutex{},
	}
}

func (repo *JsonRepository) List() ([]Student, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.readAll()
}

func (repo *JsonRepository) GetById(id int) (Student, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	students, err := repo.readAll()
	if err != nil {
		return Student{}, err
	}

	for _, s := range students {
		if s.Id == id {
			return s, nil
		}
	}
	return Student{}, ErrStudentNotFound
}

func (repo *JsonRepository) Create(s Student) (Student, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	students, err := repo.readAll()
	if err != nil {
		return Student{}, err
	}

	s.Id = generateStudentId(students)
	students = append(students, s)

	if err := repo.save(students); err != nil {
		return Student{}, err
	}
	return s, nil
}

func (repo *JsonRepository) Delete(id int) (Student, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	students, err := repo.readAll()
	if err != nil {
		return Student{}, err
	}

	targetIndex := -1
	for i, existing := range students {
		if existing.Id == id {
			targetIndex = i
			break
		}
	}
	if targetIndex == -1 {
		return Student{}, ErrStudentNotFound
	}

	deleted := students[targetIndex]
	students = slices.Delete(students, targetIndex, targetIndex+1)

	if err := repo.save(students); err != nil {
		return Student{}, err
	}
	return deleted, nil

}

func (repo *JsonRepository) Update(s Student) (Student, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	students, err := repo.readAll()
	if err != nil {
		return Student{}, err
	}

	targetIndex := -1 // impossible id
	for i, existing := range students {
		if existing.Id == s.Id {
			targetIndex = i
			break
		}
	}
	if targetIndex == -1 {
		return Student{}, ErrStudentNotFound
	}
	students[targetIndex] = s

	if err := repo.save(students); err != nil {
		return Student{}, err
	}
	return s, nil
}

func (repo *JsonRepository) readAll() ([]Student, error) {
	data, err := os.ReadFile(repo.path)
	if err != nil {
		return nil, err
	}

	var students []Student
	if err := json.Unmarshal(data, &students); err != nil {
		return nil, err
	}

	return students, nil
}

func (repo *JsonRepository) save(students []Student) error {
	data, err := json.Marshal(students)
	if err != nil {
		return err
	}
	return os.WriteFile(repo.path, data, 0644)
}

func generateStudentId(students []Student) int {
	nextId := 0
	for _, s := range students {
		if s.Id > nextId {
			nextId = s.Id
		}
	}
	return nextId + 1
}
