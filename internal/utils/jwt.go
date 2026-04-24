package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// AccessTokenDuration is how long an access token remains valid.
	AccessTokenDuration = 20 * time.Minute
	// RefreshTokenDuration is how long a refresh token remains valid.
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// TokenPair holds both the access and refresh tokens returned on sign-in.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// getAccessSecret reads the access token signing key from the environment.
func getAccessSecret() []byte {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "access-supersecret" // fallback for development only
	}
	return []byte(secret)
}

// getRefreshSecret reads the refresh token signing key from the environment.
func getRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "refresh-supersecret" // fallback for development only
	}
	return []byte(secret)
}

// GenerateAccessToken creates a short-lived JWT (20 minutes) that carries the
// user's ID and email. This token is sent with every authenticated request.
func GenerateAccessToken(userID interface{}, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userID,
		"email":  email,
		"type":   "access",
		"exp":    time.Now().Add(AccessTokenDuration).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString(getAccessSecret())
}

// GenerateRefreshToken creates a long-lived JWT (7 days) that only carries the
// user's ID. It is used solely to obtain a new access token once the old one
// expires, so it intentionally contains minimal claims.
func GenerateRefreshToken(userID interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userID,
		"type":   "refresh",
		"exp":    time.Now().Add(RefreshTokenDuration).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString(getRefreshSecret())
}

// GenerateTokenPair is a convenience function that creates both an access and
// a refresh token in a single call.
func GenerateTokenPair(userID interface{}, email string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyAccessToken parses and validates an access token JWT string.
func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	return verifyTokenWithSecret(tokenString, getAccessSecret())
}

// VerifyRefreshToken parses and validates a refresh token JWT string.
func VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return verifyTokenWithSecret(tokenString, getRefreshSecret())
}

// verifyTokenWithSecret is an internal helper that parses and validates a JWT
// against the provided secret.
func verifyTokenWithSecret(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
