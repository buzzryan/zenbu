package infra

import "github.com/buzzryan/zenbu/internal/user/usecase"

type jwsTokenManager struct{}

func (j jwsTokenManager) Generate(claims *usecase.Claims) (token string, err error) {
	panic("implement me")
}

func (j jwsTokenManager) Parse(token string) (*usecase.Claims, error) {
	panic("implement me")
}

func NewJWSTokenManager() usecase.TokenManager {
	return &jwsTokenManager{}
}
