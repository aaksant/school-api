package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	rt "github.com/aaksant/school-api/internal/api/refresh_tokens"
	"github.com/aaksant/school-api/internal/auth"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	usersRepo         UsersRepository
	refreshTokensRepo rt.RefreshTokensRepository
	validate          *validator.Validate
	tokenManager      *auth.TokenManager
	refreshTokenTtl   time.Duration
}

// dummyHash is compared against when the email doesn't exist, so that path
// costs the same as a real password check. Without it, "no such user" returns
// in ~1ms while "wrong password" takes ~50-100ms (bcrypt is slow on purpose),
// and an attacker can measure that gap to find which emails have accounts
// (a timing side channel)
//
// it must use the same cost as real password hashes, because bcrypt stores
// the cost inside the hash. A cheaper dummy would make this path faster and
// reintroduce the leak.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy"), bcrypt.DefaultCost)

func NewHandler(
	usersRepo UsersRepository,
	refreshTokensRepo rt.RefreshTokensRepository,
	validate *validator.Validate,
	tokenManager *auth.TokenManager,
	refreshTokenTtl time.Duration,
) *Handler {
	return &Handler{
		usersRepo:         usersRepo,
		refreshTokensRepo: refreshTokensRepo,
		validate:          validate,
		tokenManager:      tokenManager,
		refreshTokenTtl:   refreshTokenTtl,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

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

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("register: error hashing password: %v", err)
		http.Error(w, "hashing error", http.StatusInternalServerError)
		return
	}

	u := User{
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(passwordHash),
		Role:         auth.Role(req.Role),
		PersonID:     req.PersonID,
	}

	created, err := h.usersRepo.Create(r.Context(), u)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, ErrInvalidReference) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("register: error creating user: %v", err)
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(toUserResponse(created)); err != nil {
		log.Println("encode error:", err)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.usersRepo.GetByEmail(r.Context(), req.Email)
	// handle errors first, touch the value once you know the call succeeded
	switch {
	case errors.Is(err, ErrUserNotFound):
		// timing side channel
		// run a bcrypt compare anyway and discard the result. The call exists only
		// to spend the same time as the wrong-password path
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
		auth.Unauthorize(w)
	case err != nil:
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if !user.IsActive {
		// "user is inactive" tells an attacker the account exists and is disabled
		auth.Unauthorize(w)
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		auth.Unauthorize(w)
	}

	now := time.Now()

	accessToken, err := h.tokenManager.Sign(user.ID, user.Role, user.PersonID)
	if err != nil {
		log.Printf("login: signing token: %v", err)
		http.Error(w, "error signing access token", http.StatusInternalServerError)
		return
	}

	raw, hashed, err := rt.Generate()
	if err != nil {
		log.Printf("login: generating refresh token hash: %v", err)
		http.Error(
			w,
			"error generating refresh token hash",
			http.StatusInternalServerError,
		)
		return
	}

	// RefreshToken that Create returns is the database row (ID, hash, expiry, timestamps).
	// client only needs the raw string
	_, err = h.refreshTokensRepo.Create(r.Context(), rt.RefreshToken{
		UserId:    user.ID,
		TokenHash: hashed,
		ExpiresAt: now.Add(h.refreshTokenTtl),
	})
	if err != nil {
		log.Printf("login: creating refresh token: %v", err)
		http.Error(w, "error creating refresh token", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		AccessToken:      accessToken,
		RefreshToken:     raw,
		RefreshExpiresAt: now.Add(h.refreshTokenTtl),
		User:             toUserResponse(user),
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("encode refresh token error:", err)
	}
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now()

	hash := rt.Hash(req.RefreshToken)
	old, err := h.refreshTokensRepo.GetByHash(r.Context(), hash)
	switch {
	case errors.Is(err, rt.ErrTokenNotFound):
		auth.Unauthorize(w)
	case err != nil:
		log.Printf("refresh: get token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if old.RevokedAt != nil {
		auth.Unauthorize(w)
	}
	if now.After(old.ExpiresAt) {
		auth.Unauthorize(w)
	}

	user, err := h.usersRepo.GetById(r.Context(), old.UserId)
	if errors.Is(err, ErrUserNotFound) {
		auth.Unauthorize(w)
	}
	if !user.IsActive {
		auth.Unauthorize(w)
	}
	if err != nil {
		log.Printf("refresh: getting user %d: %v", user.ID, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	raw, hashed, err := rt.Generate()
	if err != nil {
		log.Printf("refresh: generate token: %v", err)
		http.Error(
			w,
			"error generating refresh token hash",
			http.StatusInternalServerError,
		)
		return
	}

	rotated, err := h.refreshTokensRepo.Rotate(r.Context(), old.Id, rt.RefreshToken{
		UserId:    old.UserId,
		TokenHash: hashed,
		ExpiresAt: now.Add(h.refreshTokenTtl),
	})
	if errors.Is(err, rt.ErrTokenAlreadyRevoked) {
		auth.Unauthorize(w)
	}
	if err != nil {
		log.Printf("refresh: rotate token %d: %v", old.Id, err)
		http.Error(w, "error rotating refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.tokenManager.Sign(user.ID, user.Role, user.PersonID)
	if err != nil {
		log.Printf("refresh: signing token for user %d: %v", user.ID, err)
		http.Error(w, "error signing access token", http.StatusInternalServerError)
		return
	}

	resp := RefreshResponse{
		AccessToken:      accessToken,
		RefreshToken:     raw,
		RefreshExpiresAt: rotated.ExpiresAt,
		User:             toUserResponse(user),
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("encode refresh token error:", err)
	}
}

// log out everywhere
func (h *Handler) Signout(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		log.Printf("signout: get claims: %v", claims)
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		auth.Unauthorize(w)
		return
	}

	err = h.refreshTokensRepo.RevokeAllForUser(r.Context(), userId)
	if err != nil {
		log.Printf("signout: revoke all for user: %v", err)
		http.Error(
			w,
			"error revoke all refresh token for user",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:       u.ID,
		Email:    u.Email,
		Role:     string(u.Role),
		IsActive: u.IsActive,
	}
}
