package utils

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration = 20 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}


func getAccessSecret() []byte {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "access-supersecret" // fallback for development only
	}
	return []byte(secret)
}


func getRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "refresh-supersecret" // fallback for development only
	}
	return []byte(secret)
}


func GenerateAccessToken(userID any, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userID,
		"email":  email,
		"type":   "access",
		"exp":    time.Now().Add(AccessTokenDuration).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString(getAccessSecret())
}


func GenerateRefreshToken(userID any, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userID,
		"email":  email,
		"type":   "refresh",
		"exp":    time.Now().Add(RefreshTokenDuration).Unix(),
		"iat":    time.Now().Unix(),
	})
	return token.SignedString(getRefreshSecret())
}
func GenerateTokenPair(userID any, email string) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, email)
	if err != nil {
		slog.Error("failed to generate access token", "error", err)
		return &TokenPair{}, err
	}

	refreshToken, err := GenerateRefreshToken(userID, email)
	if err != nil {
		slog.Error("failed to generate refresh token", "error", err)
		return &TokenPair{}, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	return verifyTokenWithSecret(tokenString, getAccessSecret())
}

func VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return verifyTokenWithSecret(tokenString, getRefreshSecret())
}

func verifyTokenWithSecret(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("unexpected signing method", "method", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		slog.Error("invalid token", "error", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Error("invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
