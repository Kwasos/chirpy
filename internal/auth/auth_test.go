package auth

import "testing"

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
