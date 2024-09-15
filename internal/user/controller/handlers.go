package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
	"github.com/buzzryan/zenbu/internal/user/domain"
	"github.com/buzzryan/zenbu/internal/user/usecase"
)

// BasicSignupCtrl is a controller for basic signup.
type BasicSignupCtrl struct {
	uc usecase.BasicSignupUC
}

type BasicSignupReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewBasicSignupCtrl(uc usecase.BasicSignupUC) *BasicSignupCtrl {
	return &BasicSignupCtrl{uc: uc}
}

const (
	CodeUsernameAlreadyExists = 2000
)

func (b *BasicSignupCtrl) Handle(w http.ResponseWriter, req *http.Request) error {
	var reqBody BasicSignupReq
	err := httputil.ParseJSONBody(req, &reqBody)
	if err != nil {
		return httputil.HandleParseJSONBodyError(req.Context(), w, err)
	}

	logutil.From(req.Context()).Info("Basic Signup", slog.Any("req", reqBody))

	// fixme: use validator to reduce boilerplate
	if reqBody.Username == "" || reqBody.Password == "" {
		return httputil.ResponseError(w, http.StatusBadRequest, 0, "username and password is required")
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
