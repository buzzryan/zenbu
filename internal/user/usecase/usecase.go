package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/commonutil/storageutil"
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
	userRepo     UserRepo
	tokenManager TokenManager
}

func NewBasicSignupUC(userRepo UserRepo, manager TokenManager) BasicSignupUC {
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
	userRepo UserRepo
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

func NewAuthenticateUC(userRepo UserRepo, manager TokenManager) AuthenticateUC {
	return &authenticateUC{userRepo: userRepo, manager: manager}
}

type BasicLoginRes struct {
	Token string
}

type BasicLoginUC interface {
	Execute(ctx context.Context, username, password string) (*BasicLoginRes, error)
}

type basicLoginUC struct {
	userRepo     UserRepo
	tokenManager TokenManager
}

func (b basicLoginUC) Execute(ctx context.Context, username, password string) (*BasicLoginRes, error) {
	u, err := b.userRepo.GetByName(ctx, username)
	if err != nil {
		return nil, err
	}

	if !u.Password.Compare(password) {
		return nil, ErrInvalidPassword
	}

	token, err := b.tokenManager.Generate(&Claims{
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(TokenExpiresIn),
	})
	if err != nil {
		return nil, err
	}

	return &BasicLoginRes{Token: token}, nil
}

func NewBasicLoginUC(userRepo UserRepo, manager TokenManager) BasicLoginUC {
	return &basicLoginUC{userRepo: userRepo, tokenManager: manager}
}

// CreateProfileImagUploadURLUC returns a signed URL for uploading a profile image.
type CreateProfileImagUploadURLUC interface {
	Execute(ctx context.Context, token string) (url string, err error)
}

type createProfileImageUploadURL struct {
	userRepo     UserRepo
	tokenManager TokenManager
	storage      storageutil.Storage
}

func NewCreateProfileImagUploadURLUC(userRepo UserRepo, tokenManager TokenManager, storage storageutil.Storage) CreateProfileImagUploadURLUC {
	return &createProfileImageUploadURL{userRepo: userRepo, tokenManager: tokenManager, storage: storage}
}

func userProfileImageDir(userID uuid.UUID) string {
	return "profiles/" + userID.String() + "/images"
}

func (c *createProfileImageUploadURL) Execute(ctx context.Context, token string) (string, error) {
	claims, err := c.tokenManager.Parse(token)
	if err != nil {
		return "", err
	}

	url, err := c.storage.CreateUploadURL(ctx, storageutil.Public,
		userProfileImageDir(claims.UserID)+"/"+time.Now().Format("20060102150405.999999999"))
	if err != nil {
		return "", err
	}

	return url, nil
}

type GetProfileImageURLUC interface {
	Execute(ctx context.Context, userID uuid.UUID) (url string, err error)
}

type getProfileImageURLUC struct {
	userRepo UserRepo
	storage  storageutil.Storage
}

func NewGetProfileImageURLUC(userRepo UserRepo, storage storageutil.Storage) GetProfileImageURLUC {
	return &getProfileImageURLUC{userRepo: userRepo, storage: storage}
}

func (g *getProfileImageURLUC) Execute(ctx context.Context, userID uuid.UUID) (string, error) {
	files, err := g.storage.ListFiles(ctx, storageutil.Public, userProfileImageDir(userID))
	if err != nil {
		return "", fmt.Errorf("failed to get image url: %w", err)
	}

	if len(files) == 0 {
		return "", errors.New("image not found")
	}

	slices.SortFunc(files, func(i, j *storageutil.File) int {
		return j.UpdatedAt.Compare(i.UpdatedAt)
	})

	return g.storage.GetPublicFileURL(ctx, files[0].Filepath)
}
