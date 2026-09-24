package teachingassignments

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
	repo     TeachingAssignmentsRepository
	validate *validator.Validate
}

func NewHandler(
	repo TeachingAssignmentsRepository,
	validate *validator.Validate,
) *Handler {
	return &Handler{
		repo:     repo,
		validate: validate,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	teacherId, err := strconv.Atoi(r.PathValue("teacherId"))
	if err != nil {
		http.Error(w, "invalid teacher ID", http.StatusBadRequest)
		return
	}

	mappings, err := h.repo.List(r.Context(), teacherId)
	if err != nil {
		http.Error(w, "error listing teacher mappings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mappings); err != nil {
		log.Println("encode error:", err)
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid mapping id", http.StatusBadRequest)
	}

	mapping, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTeachingAssignmentNotFound) {
			http.Error(w, "no rows retrieved", http.StatusNoContent)
			return
		}
		http.Error(w, "error fetching mapping", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapping); err != nil {
		log.Println("encode error:", err)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	teacherId, err := strconv.Atoi(r.PathValue("teacherId"))
	if err != nil {
		http.Error(w, "invalid teacher ID", http.StatusBadRequest)
		return
	}

	var req CreateTeachingAssignmentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mapping := TeachingAssignment{
		TeacherId: teacherId,
		ClassId:   req.ClassId,
		SubjectId: req.SubjectId,
	}

	created, err := h.repo.Create(r.Context(), mapping)
	if err != nil {
		if errors.Is(err, ErrInvalidReference) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrDuplicate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "error creating mapping", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Print("encode error:", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid mapping id", http.StatusBadRequest)
		return
	}

	mapping, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "mapping not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "error deleting mapping", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mapping); err != nil {
		log.Println("encode error:", err)
	}
}
