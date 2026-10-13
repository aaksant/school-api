package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleStudent   Role = "student"
	RoleTeacher   Role = "teacher"
	RoleStaff     Role = "staff"
	RoleExecutive Role = "executive"
)

type Claims struct {
	jwt.RegisteredClaims
	Role     Role  `json:"role"`
	PersonId int64 `json:"person_id"`
}

type ctxKey int

const claimsKey ctxKey = 0

// 32 bytes/256 bits is considered enough to brute force sha-256
const TokenBytes = 32

type TokenManager struct {
	key       []byte
	accessTtl time.Duration
	issuer    string
}

func NewTokenManager(
	key []byte,
	accessTtl time.Duration,
	issuer string,
) (*TokenManager, error) {
	if len(key) < TokenBytes {
		return nil, fmt.Errorf("key must be at least 32 bytes: %v", len(key))
	}
	if accessTtl <= 0 {
		return nil, fmt.Errorf("ttl must be positive")
	}

	return &TokenManager{
		key:       key,
		accessTtl: accessTtl,
		issuer:    issuer,
	}, nil
}

func (tm *TokenManager) Sign(
	userId int64,
	role Role,
	personId int64,
) (string, error) {
	if !role.Valid() {
		return "", fmt.Errorf("role does not exist")
	}

	now := time.Now()
	claims := Claims{
		Role:     role,
		PersonId: personId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userId, 10),
			Issuer:    tm.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.accessTtl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.key)
}

func (tm *TokenManager) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return tm.key, nil
		},
		jwt.WithValidMethods([]string{"HS256"}), // skip alg check in keyFunc, enforce HS256
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(tm.issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims) // assert *Claims
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	if !claims.Role.Valid() {
		return nil, fmt.Errorf("invalid role on claims")
	}

	return claims, nil
}

func (r Role) Valid() bool {
	switch r {
	case RoleStudent, RoleTeacher, RoleStaff, RoleExecutive:
		return true
	}
	return false
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}
