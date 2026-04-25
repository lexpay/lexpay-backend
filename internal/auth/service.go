package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/luponetn/lexpay/internal/db"
	"github.com/luponetn/lexpay/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo db.Querier
}

type Svc interface {
	SignUp(ctx context.Context, arg db.CreateUserOnSignupParams) (db.Users, error)
	SignIn(ctx context.Context, email, password string) (*utils.TokenPair, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*utils.TokenPair, error)
}

func NewService(repo db.Querier) Svc {
	return &Service{repo: repo}
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

func (s *Service) SignUp(ctx context.Context, arg db.CreateUserOnSignupParams) (db.Users, error) {
	arg.Email = strings.ToLower(strings.TrimSpace(arg.Email))
	_, err := s.repo.FindUserByEmail(ctx, arg.Email)
	if err == nil {
		slog.Error("user already exists", "email", arg.Email)
		return db.Users{}, ErrUserAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("error finding user", "error", err)
		return db.Users{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(arg.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return db.Users{}, err
	}
	arg.PasswordHash = string(hashedPassword)

	return s.repo.CreateUserOnSignup(ctx, arg)
}

func (s *Service) SignIn(ctx context.Context, email, password string) (*utils.TokenPair, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		slog.Error("error finding user during signin", "error", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	tokens, err := utils.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		slog.Error("failed to generate token pair", "error", err)
		return nil, err
	}

	return tokens, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshTokenString string) (*utils.TokenPair, error) {
	claims, err := utils.VerifyRefreshToken(refreshTokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["userid"]
	if !ok {
		return nil, ErrInvalidToken
	}

	email, _ := claims["email"].(string)
	if email == "" {
		return nil, ErrInvalidToken
	}

	tokens, err := utils.GenerateTokenPair(userID, email)
	if err != nil {
		slog.Error("failed to generate token pair on refresh", "error", err)
		return nil, err
	}

	return tokens, nil
}

