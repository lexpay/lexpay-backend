package auth

import (
	"context"

	"github.com/luponetn/lexpay/internal/db"
)

type Service struct {
	repo db.Querier
}

type Svc interface {
	SignUp(ctx context.Context, arg db.CreateUserOnSignupParams) (db.Users, error)
}

func NewService(repo db.Querier) Svc {
	return &Service{repo: repo}
}


//implement auth services
func (s *Service) SignUp(ctx context.Context, arg db.CreateUserOnSignupParams) (db.Users, error) {
	return s.repo.CreateUserOnSignup(ctx, arg)
}
