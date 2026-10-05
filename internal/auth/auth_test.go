package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestHashAndCheck(t *testing.T) {
	hash, err := HashPassword("correct")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	ok, err := CheckPasswordHash("correct", hash)
	if err != nil || !ok {
		t.Errorf("expected match, got ok=%v err=%v", ok, err)
	}
}

func TestCheckPasswordHashWrong(t *testing.T) {
	hash, err := HashPassword("correct")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	ok, err := CheckPasswordHash("wrong", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected mismatch, got match")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotID != userID {
		t.Errorf("expected %v, got %v", userID, gotID)
	}
}

func TestValidateWrongSecretJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	secret2 := "secret-test"

	token, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ValidateJWT(token, secret2)
	if err == nil {
		t.Fatalf("unexpected lack of error for wrong secret")
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer some-token-value")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "some-token-value" {
		t.Errorf("expected 'some-token-value', got '%s'", token)
	}
}

func TestGetBearerBadToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Auth", "Bearer some-token-value")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error `missing bearer token`, got nil")
	}

	if err.Error() != "missing bearer token" {
		t.Fatalf("expected error `missing bearer token`, got: %v", err)
	}
}

func TestGetBearerBadBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bear some-token-value")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected error `bearer token not found`, got nil")
	}

	if err.Error() != "bearer token not found" {
		t.Fatalf("expected error `bearer token not found`, got: %v", err)
	}
}
