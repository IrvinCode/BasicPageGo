package main

import (
	"net/http"
	"github.com/bmizerany/pat"
)

func (app *App) Routes() http.Handler {
	// Wrap all of our web page route with the NoSurf middleware.
	mux := pat.New()
	mux.Get("/", NoSurf(app.Home))
	mux.Get("/snippet/new", app.RequireLogin(NoSurf(app.NewSnippet)))
	mux.Post("/snippet/new", app.RequireLogin(NoSurf(app.CreateSnippet)))
	mux.Get("/snippet/delete", app.RequireAdmin(NoSurf(app.EraseSnippet)))
	mux.Post("/snippet/delete", app.RequireAdmin(NoSurf(app.DeleteSnippet)))
	//mux.Get("/snippet/:id", NoSurf(app.ShowSnippet))
	mux.Get("/snippet/:id", app.RequireLogin(NoSurf(app.ShowSnippet)))

	mux.Get("/user/signup", NoSurf(app.SignupUser))
	mux.Post("/user/signup", NoSurf(app.CreateUser))
	mux.Get("/user/login", NoSurf(app.LoginUser))
	mux.Post("/user/login", NoSurf(app.VerifyUser))
	mux.Post("/user/logout", app.RequireLogin(NoSurf(app.LogoutUser)))

	mux.Get("/admin/signup", app.RequireAdmin(NoSurf(app.SignupAdmin)))
	mux.Post("/admin/signup", app.RequireAdmin(NoSurf(app.CreateAdmin)))

	fileServer := http.FileServer(http.Dir(app.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return LogRequest(SecureHeaders(mux))
}
