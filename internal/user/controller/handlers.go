package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
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

type BasicSignupRes struct {
	Token string `json:"token"`
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
	if errors.Is(err, usecase.ErrUsernameAlreadyExists) {
		return httputil.ResponseError(w, http.StatusConflict, CodeUsernameAlreadyExists, "username already exists")
	}
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute Basic Signup", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, &BasicSignupRes{Token: res.Token})
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
	if errors.Is(err, usecase.ErrUserNotFound) {
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

type BasicLoginCtrl struct {
	uc usecase.BasicLoginUC
}

func NewBasicLoginCtrl(uc usecase.BasicLoginUC) *BasicLoginCtrl {
	return &BasicLoginCtrl{uc: uc}
}

type BasicLoginReq struct {
	Username string `json:"username" validate:"required,max=32,min=1"`
	Password string `json:"password" validate:"required,password"`
}

type BasicLoginRes struct {
	Token string `json:"token"`
}

func (b *BasicLoginCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	var reqBody BasicLoginReq
	if err := httputil.ParseJSONBody(req, &reqBody); err != nil {
		return httputil.HandleParseJSONBodyError(req.Context(), w, err)
	}

	if err := validutil.Validate(reqBody); err != nil {
		return httputil.ResponseError(w, http.StatusBadRequest, httputil.CodeInvalidRequestParams, err.Error())
	}

	res, err := b.uc.Execute(req.Context(), reqBody.Username, reqBody.Password)
	if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrInvalidPassword) {
		return httputil.ResponseError(w, http.StatusUnauthorized, httputil.CodeUnauthenticated, "invalid credentials")
	}
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute Basic Signup", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, &BasicLoginRes{Token: res.Token})
}

type CreateProfileImageUploadURLCtrl struct {
	uc usecase.CreateProfileImagUploadURLUC
}

type CreateProfileImageUploadURLRes struct {
	URL string `json:"url"`
}

func NewCreateProfileImageUploadURLCtrl(uc usecase.CreateProfileImagUploadURLUC) *CreateProfileImageUploadURLCtrl {
	return &CreateProfileImageUploadURLCtrl{uc: uc}
}

func (c *CreateProfileImageUploadURLCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	token, err := httputil.GetBearerToken(req)
	if err != nil {
		return httputil.ResponseError(w, http.StatusUnauthorized, httputil.CodeUnauthenticated, err.Error())
	}

	url, err := c.uc.Execute(req.Context(), token)
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute CreateProfileImagUploadURL", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, &CreateProfileImageUploadURLRes{URL: url})
}

type GetProfileImageURLCtrl struct {
	uc usecase.GetMyProfileImageURLUC
}

func NewGetProfileImageURLCtrl(uc usecase.GetMyProfileImageURLUC) *GetProfileImageURLCtrl {
	return &GetProfileImageURLCtrl{uc: uc}
}

func (g *GetProfileImageURLCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	userID := req.PathValue("id")
	if userID == "" {
		return httputil.ResponseError(w, http.StatusBadRequest, httputil.CodeInvalidRequestParams, "user id required")
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return httputil.ResponseError(w, http.StatusBadRequest, httputil.CodeInvalidRequestParams, "invalid user id")
	}

	url, err := g.uc.Execute(req.Context(), parsedUserID)
	if err != nil {
		logutil.From(req.Context()).Error("failed to execute GetProfileImageURL", slog.Any("err", err))
		return httputil.ResponseError(w, http.StatusInternalServerError, 0, "internal server error")
	}

	return httputil.ResponseJSON(w, http.StatusOK, &CreateProfileImageUploadURLRes{URL: url})
}
