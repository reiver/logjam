package tokens

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

// Initialize inputSecret is optional. If nil or shorter than 32 bytes, a random secret is generated.
func Initialize(inputSecret []byte) {
	if inputSecret != nil && len(inputSecret) >= 32 {
		secret = inputSecret
		return
	}

	secret = make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		log.Fatalf("failed to generate secret: %v", err)
	}
	log.Printf("JWT secret (hex): %s", hex.EncodeToString(secret))
}

// CreateToken creates a JWT for a user ID
func CreateToken(userID string, extraData map[string]any) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(32 * 24 * time.Hour).Unix(),
	}
	if extraData != nil && len(extraData) > 0 {
		for k, v := range extraData {
			claims[k] = v
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken decodes and validates a JWT
func ParseToken(tokenStr string) (map[string]any, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userId"] == nil {
		return claims, errors.New("invalid claims")
	}

	return claims, nil
}
