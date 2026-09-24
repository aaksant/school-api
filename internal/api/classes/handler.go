package classes

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
	repo     ClassesRepository
	validate *validator.Validate
}

func NewHandler(
	repo ClassesRepository,
	validate *validator.Validate,
) *Handler {
	return &Handler{
		repo:     repo,
		validate: validate,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	classes, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "error listing classes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(classes); err != nil {
		log.Println("classes.List: encode error:", err)
	}
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid class ID", http.StatusBadRequest)
		return
	}

	class, err := h.repo.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching class", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(class); err != nil {
		log.Println("classes.GetById: encode error:", err)
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateClassRequest

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

	c := Class{
		Name:              req.Name,
		HomeroomTeacherId: req.HomeroomTeacherId,
	}

	created, err := h.repo.Create(r.Context(), c)
	if err != nil {
		if errors.Is(err, ErrInvalidReference) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, ErrDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "error creating class", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Println("encode error:", err)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid class ID", http.StatusBadRequest)
		return
	}

	class, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "class not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "error deleting class", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(class); err != nil {
		log.Println("encode error:", err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid class ID", http.StatusBadRequest)
		return
	}

	var req UpdateClassRequest

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

	c := Class{
		Id:                id,
		Name:              req.Name,
		HomeroomTeacherId: req.HomeroomTeacherId,
	}

	updated, err := h.repo.Update(r.Context(), c)
	if err != nil {
		if errors.Is(err, ErrClassNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "error updating class", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Println("encode error:", err)
	}
}
