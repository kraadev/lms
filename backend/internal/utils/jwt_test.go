package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	password := "SecretP@ssw0rd!2026"

	// 1. Test successful hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == password {
		t.Errorf("expected hash to differ from raw password")
	}

	// 2. Test valid verification
	if !CheckPasswordHash(password, hash) {
		t.Errorf("expected password verification to succeed for correct password")
	}

	// 3. Test invalid verification
	if CheckPasswordHash("WrongPassword", hash) {
		t.Errorf("expected password verification to fail for incorrect password")
	}
}

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "super-secure-jwt-secret-key-32-chars-long!"
	userID := int64(42)
	email := "alex@lms.local"
	role := "teacher"
	name := "Alex Teacher"
	expiryHours := 2

	// 1. Generate token
	tokenStr, err := GenerateJWT(userID, email, role, name, secret, expiryHours)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}

	// 2. Validate token
	claims, err := ValidateJWT(tokenStr, secret)
	if err != nil {
		t.Fatalf("failed to validate valid JWT: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected Role %s, got %s", role, claims.Role)
	}
	if claims.Name != name {
		t.Errorf("expected Name %s, got %s", name, claims.Name)
	}
}

func TestValidateJWT_InvalidSecret(t *testing.T) {
	validSecret := "primary-secret-key-12345678901234"
	wrongSecret := "different-secret-key-123456789012"

	tokenStr, err := GenerateJWT(1, "test@lms.local", "student", "Student", validSecret, 1)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	_, err = ValidateJWT(tokenStr, wrongSecret)
	if err == nil {
		t.Errorf("expected validation to fail when using wrong secret, but succeeded")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	secret := "secret-for-expired-token-test-123"

	// Create an already expired token manually
	claims := JWTClaims{
		UserID: 99,
		Email:  "expired@lms.local",
		Role:   "student",
		Name:   "Expired User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "lms-backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = ValidateJWT(tokenStr, secret)
	if err == nil {
		t.Errorf("expected validation to fail for expired token, but succeeded")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	secret := "valid-secret-1234567890"

	malformedTokens := []string{
		"",
		"not-a-valid-jwt",
		"header.payload",
		"header.payload.signature.extra",
	}

	for _, token := range malformedTokens {
		_, err := ValidateJWT(token, secret)
		if err == nil {
			t.Errorf("expected validation to fail for malformed token '%s'", token)
		}
	}
}
