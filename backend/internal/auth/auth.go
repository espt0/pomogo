package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GerarAccessToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GerarRefreshToken() (raw, hash string, err error) {
	b := make([]byte, 32)

	// Gera (token) números aleatórios criptograficamente seguros
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	// Transforma em string e manda pro cliente
	raw = hex.EncodeToString(b)

	// Calcula o hash do token
	h := sha256.Sum256([]byte(raw))

	// Transforma em string e manda pro banco
	hash = hex.EncodeToString(h[:])

	return raw, hash, nil
}

func HashRefreshToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
