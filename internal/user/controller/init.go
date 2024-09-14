package controller

import (
	"net/http"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/user/domain"
	"github.com/buzzryan/zenbu/internal/user/usecase"
)

type InitOpts struct {
	Mux          *http.ServeMux
	UserRepo     domain.UserRepo
	TokenManager usecase.TokenManager
}

func Init(opts *InitOpts) {
	basicSignupUC := usecase.NewBasicSignupUC(opts.UserRepo, opts.TokenManager)
	basicSignupCtrl := NewBasicSignupCtrl(basicSignupUC)

	// register routers
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/signup", basicSignupCtrl.Handle)
}
