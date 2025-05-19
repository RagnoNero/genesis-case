package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type TokenGenerator struct {
	saltLength int
}

func NewTokenGenerator(saltLength int) *TokenGenerator {
	return &TokenGenerator{
		saltLength: saltLength,
	}
}

func (t *TokenGenerator) Generate(email string) (string, error) {
	salt, err := generateRandomSalt(t.saltLength)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	raw := fmt.Sprintf("%s:%s:%s", email, timestamp, salt)

	hash := sha256.Sum256([]byte(raw))
	token := base64.URLEncoding.EncodeToString(hash[:])

	return token, nil
}

func generateRandomSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
