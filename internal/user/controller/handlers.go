package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
	"github.com/buzzryan/zenbu/internal/user/domain"
	"github.com/buzzryan/zenbu/internal/user/usecase"
	"github.com/buzzryan/zenbu/internal/validutil"
)

const (
	CodeUsernameAlreadyExists = 2000
	CodeUserNotFound          = 2001
)

// BasicSignupCtrl is a controller for basic signup.
type BasicSignupCtrl struct {
	uc usecase.BasicSignupUC
}

type BasicSignupReq struct {
	Username string `json:"username" validate:"required,max=32,min=1"`
	Password string `json:"password" validate:"required,password"`
}

func NewBasicSignupCtrl(uc usecase.BasicSignupUC) *BasicSignupCtrl {
	return &BasicSignupCtrl{uc: uc}
}

func (b *BasicSignupCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	var reqBody BasicSignupReq
	if err := httputil.ParseJSONBody(req, &reqBody); err != nil {
		return httputil.HandleParseJSONBodyError(req.Context(), w, err)
	}

	if err := validutil.Validate(reqBody); err != nil {
		return httputil.ResponseError(w, http.StatusBadRequest, httputil.CodeInvalidRequestParams, err.Error())
	}

	res, err := b.uc.Execute(req.Context(), &usecase.SignupReq{
		Username: reqBody.Username,
		Password: reqBody.Password,
	})
	if errors.Is(err, domain.ErrUsernameAlreadyExists) {
		return httputil.ResponseError(w, http.StatusConflict, CodeUsernameAlreadyExists, "username already exists")
	}
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute Basic Signup", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, res)
}

type AuthenticateCtrl struct {
	uc usecase.AuthenticateUC
}

func NewAuthenticateCtrl(uc usecase.AuthenticateUC) *AuthenticateCtrl {
	return &AuthenticateCtrl{uc: uc}
}

type AuthenticateRes struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

func (a *AuthenticateCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	token, err := httputil.GetBearerToken(req)
	if err != nil {
		return httputil.ResponseError(w, http.StatusUnauthorized, httputil.CodeUnauthenticated, err.Error())
	}

	res, err := a.uc.Execute(req.Context(), token)
	if errors.Is(err, usecase.ErrInvalidToken) {
		return httputil.ResponseError(w, http.StatusUnauthorized, httputil.CodeUnauthenticated, err.Error())
	}
	if errors.Is(err, usecase.ErrTokenExpired) {
		return httputil.ResponseError(w, http.StatusUnauthorized, httputil.CodeTokenExpired, err.Error())
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return httputil.ResponseError(w, http.StatusNotFound, CodeUserNotFound, err.Error())
	}
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute Authenticate", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, &AuthenticateRes{
		Token:    res.RefreshedToken,
		UserID:   res.User.ID.String(),
		Username: res.User.Username,
	})
}
