package controller

import (
	"net/http"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/user/usecase"
)

type InitOpts struct {
	Mux          *http.ServeMux
	UserRepo     usecase.UserRepo
	TokenManager usecase.TokenManager
}

func Init(opts *InitOpts) {
	basicSignupUC := usecase.NewBasicSignupUC(opts.UserRepo, opts.TokenManager)
	basicSignupCtrl := NewBasicSignupCtrl(basicSignupUC)

	authenticateUC := usecase.NewAuthenticateUC(opts.UserRepo, opts.TokenManager)
	authenticateCtrl := NewAuthenticateCtrl(authenticateUC)

	basicLoginUC := usecase.NewBasicLoginUC(opts.UserRepo, opts.TokenManager)
	basicLoginCtrl := NewBasicLoginCtrl(basicLoginUC)

	// register routers
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/signup", basicSignupCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/authenticate", authenticateCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/login", basicLoginCtrl.Handle)
}
