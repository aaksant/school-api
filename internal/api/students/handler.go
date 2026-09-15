package students

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	repo     StudentsRepository
	validate *validator.Validate
}

func NewHandler(
	repo StudentsRepository,
	validate *validator.Validate,
) *Handler {
	return &Handler{
		repo:     repo,
		validate: validate,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Error listing students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(students); err != nil {
		log.Printf("students.List: encode error: %v", err) // already OK above
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No rows retrieved", http.StatusNoContent)
		}
		http.Error(w, "Error fetching student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(student); err != nil {
		log.Printf("students.GetById: encode error: %v\n", err)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateStudentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := Student{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		DateOfBirth: req.DateOfBirth,
	}

	created, err := h.repo.Create(r.Context(), s)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No rows retrieved", http.StatusNoContent)
			return
		}
		http.Error(w, "Error creating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("students.Create: encode error: %v\n", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(student); err != nil {
		log.Printf("encode error: %v\n", err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var req UpdateStudentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := Student{
		Id:          id,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		DateOfBirth: req.DateOfBirth,
	}

	updated, err := h.repo.Update(r.Context(), s)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error updating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("students.Update: encode error: %v\n", err)
	}
}

func (h *Handler) PartialUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	patched, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error patching student", http.StatusInternalServerError)
		return
	}

	var req PatchStudentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.FirstName != nil {
		patched.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		patched.LastName = *req.LastName
	}
	if req.Email != nil {
		patched.Email = *req.Email
	}

	updated, err := h.repo.Update(r.Context(), patched)
	if err != nil {
		if errors.Is(err, ErrStudentNotFound) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error patching student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("students.Update: encode error: %v\n", err)
	}
}
