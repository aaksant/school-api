package refresh_tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/aaksant/school-api/internal/auth"
)

func Generate() (raw, hashed string, err error) {
	b := make([]byte, auth.TokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	raw = base64.RawStdEncoding.EncodeToString(b)
	return raw, Hash(raw), nil
}

func Hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
