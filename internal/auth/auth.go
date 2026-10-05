// Package auth provides helper functions for authentication actions
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expiry := time.Now().UTC().Add(time.Hour)

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiry),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	stringID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	userID, err := uuid.Parse(stringID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("missing bearer token")
	}
	strippedToken, ok := strings.CutPrefix(bearerToken, "Bearer ")
	if !ok {
		return "", errors.New("bearer token not found")
	}
	return strippedToken, nil
}

func MakeRefreshToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	hexString := hex.EncodeToString(b)
	return hexString
}
