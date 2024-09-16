package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/user/domain"
)

const (
	TokenExpiresIn = time.Hour * 24 * 28 // 4 weeks
)

type SignupReq struct {
	Username string
	Password string
}

type SignupRes struct {
	Token string
}

type BasicSignupUC interface {
	Execute(ctx context.Context, req *SignupReq) (*SignupRes, error)
}

type basicSignupUC struct {
	userRepo     domain.UserRepo
	tokenManager TokenManager
}

func NewBasicSignupUC(userRepo domain.UserRepo, manager TokenManager) BasicSignupUC {
	return &basicSignupUC{userRepo: userRepo, tokenManager: manager}
}

func (b *basicSignupUC) Execute(ctx context.Context, req *SignupReq) (*SignupRes, error) {
	newUser := &domain.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Password:  domain.NewPassword(req.Password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	newUser, err := b.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	token, err := b.tokenManager.Generate(&Claims{
		UserID:    newUser.ID,
		ExpiresAt: time.Now().Add(TokenExpiresIn),
	})
	if err != nil {
		return nil, err
	}

	return &SignupRes{Token: token}, nil
}

type AuthenticateRes struct {
	User *domain.User

	RefreshedToken string
}

type AuthenticateUC interface {
	Execute(ctx context.Context, token string) (*AuthenticateRes, error)
}

type authenticateUC struct {
	userRepo domain.UserRepo
	manager  TokenManager
}

func (a authenticateUC) Execute(ctx context.Context, token string) (*AuthenticateRes, error) {
	claims, err := a.manager.Parse(token)
	if err != nil {
		return nil, err
	}

	u, err := a.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	refreshedToken, err := a.manager.Generate(&Claims{
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(TokenExpiresIn),
	})
	if err != nil {
		return nil, err
	}

	return &AuthenticateRes{User: u, RefreshedToken: refreshedToken}, nil
}

func NewAuthenticateUC(userRepo domain.UserRepo, manager TokenManager) AuthenticateUC {
	return &authenticateUC{userRepo: userRepo, manager: manager}
}
