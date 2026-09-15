package teachers

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
	repo     TeachersRepository
	validate *validator.Validate
}

func NewHandler(
	repo TeachersRepository,
	validate *validator.Validate,
) *Handler {
	return &Handler{
		repo:     repo,
		validate: validate,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "error listing teachers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(teachers); err != nil {
		log.Printf("teachers.List: encode error: %v", err) // already OK above
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid teacher ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(teacher); err != nil {
		log.Printf("teachers.GetById: encode error: %v\n", err)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTeacherRequest
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

	t := Teacher{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		DateOfBirth: req.DateOfBirth,
	}

	created, err := h.repo.Create(r.Context(), t)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no rows retrieved", http.StatusNotFound)
			return
		}
		http.Error(w, "error creating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("teachers.Create: encode error: %v\n", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid student ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(teacher); err != nil {
		log.Printf("encode error: %v\n", err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid teacher ID", http.StatusBadRequest)
		return
	}

	var t Teacher

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.Id = id

	if err := h.validate.Struct(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.repo.Update(r.Context(), t)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			http.Error(w, "teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("encode error: %v\n", err)
	}
}
