package controller

import (
	"net/http"

	"github.com/buzzryan/zenbu/internal/commonutil/httputil"
	"github.com/buzzryan/zenbu/internal/commonutil/storageutil"
	"github.com/buzzryan/zenbu/internal/user/usecase"
)

type InitOpts struct {
	Mux          *http.ServeMux
	UserRepo     usecase.UserRepo
	TokenManager usecase.TokenManager
	Storage      storageutil.Storage
}

func Init(opts *InitOpts) {
	basicSignupUC := usecase.NewBasicSignupUC(opts.UserRepo, opts.TokenManager)
	basicSignupCtrl := NewBasicSignupCtrl(basicSignupUC)

	authenticateUC := usecase.NewAuthenticateUC(opts.UserRepo, opts.TokenManager)
	authenticateCtrl := NewAuthenticateCtrl(authenticateUC)

	basicLoginUC := usecase.NewBasicLoginUC(opts.UserRepo, opts.TokenManager)
	basicLoginCtrl := NewBasicLoginCtrl(basicLoginUC)

	createProfileImageUploadUC := usecase.NewCreateProfileImagUploadURLUC(opts.UserRepo, opts.TokenManager, opts.Storage)
	createProfileImageUploadURLCtrl := NewCreateProfileImageUploadURLCtrl(createProfileImageUploadUC)

	getProfileImageURLUC := usecase.NewGetProfileImageURLUC(opts.UserRepo, opts.Storage)
	getProfileImageURLCtrl := NewGetProfileImageURLCtrl(getProfileImageURLUC)

	getMeUC := usecase.NewGetMeUC(opts.UserRepo, opts.TokenManager)
	getMeCtrl := NewGetMeCtrl(getMeUC)

	updateMyProfileUC := usecase.NewUpdateMyProfileUC(opts.UserRepo, opts.TokenManager)
	updateMyProfileCtrl := NewUpdateMyProfileCtrl(updateMyProfileUC)

	// register routers
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/signup", basicSignupCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/authenticate", authenticateCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/login", basicLoginCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPost, "/me/profile/image", createProfileImageUploadURLCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodGet, "/users/{id}/profile/image", getProfileImageURLCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodGet, "/me", getMeCtrl.Handle)
	httputil.RegisterHandler(opts.Mux, http.MethodPatch, "/me/profile", updateMyProfileCtrl.Handle)
}
