package main

import (
	"net/http"
	"strconv"
	"snippetbox.org/pkg/forms"
	"snippetbox.org/pkg/models"
	"fmt"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.Database.LatestSnippets()
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Snippets: snippets,
	})
}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.Database.GetSnippet(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if snippet == nil {
		app.NotFound(w)
		return
	}

	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "show.page.html", &HTMLData{
	Flash:   flash,
	Snippet: snippet,
	})
}

func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "new.page.html", &HTMLData{
		Form: &forms.NewSnippet{},
	})
}
func (app *App) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.NewSnippet{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: r.PostForm.Get("expires"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "new.page.html", &HTMLData{Form: form})
		return
	}

	session := app.Sessions.Load(r)

	err = session.PutString(w, "flash", "Your snippet was saved successfully!")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	id, err := app.Database.InsertSnippet(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *App) DeleteSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.DeleteSnippet{
		Id: r.PostForm.Get("id"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "delete.page.html", &HTMLData{Form: form})
		return
	}

	session := app.Sessions.Load(r)

	err = session.PutString(w, "flash", "Your snippet was deleted successfully!")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = app.Database.DeleteSnippet(form.Id)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *App) SignupUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.page.html", &HTMLData{
		Form: &forms.SignupUser{},
	})
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form := &forms.SignupUser{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	if !form.Valid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	}
	// Try to create a new user record in the database. If the email already exists
	// add a failure message to the form and re-display the form.
	err = app.Database.InsertUser(form.Name, form.Email, form.Password)
	if err == models.ErrDuplicateEmail {
		form.Failures["Email"] = "Address is already in use"
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}
	// Otherwise, add a confirmation flash message to the session confirming that
	// their signup worked and asking them to log in.
	msg := "Your signup was successful. Please log in using your credentials."
	session := app.Sessions.Load(r)
	err = session.PutString(w, "flash", msg)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// And redirect the user to the login page.
	http.Redirect(w, r, "/	", http.StatusSeeOther)
}

func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "login.page.html", &HTMLData{
		Flash: flash,
		Form: &forms.LoginUser{},
	})
}

func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.LoginUser{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map, and re-display the login page.
	currentUserID, err, admin := app.Database.VerifyUser(form.Email, form.Password)
	if err == models.ErrInvalidCredentials {
		form.Failures["Generic"] = "Email or Password is incorrect"
		app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	if admin == false {
		// Add the ID of the current user to the session, so that they are now 'logged
		// in'.
		session := app.Sessions.Load(r)
		err = session.PutInt(w, "currentUserID", currentUserID)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Redirect the user to the Add Snippet page.
		http.Redirect(w, r, "/snippet/new", http.StatusSeeOther)
	}

	// Add the ID of the current admin user to the session, so that they are now 'logged
	// in'.
	session := app.Admin.Load(r)
	err = session.PutInt(w, "currentAdminID", currentUserID)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	if err == nil {
		// Add the ID of the current user to the session, so that they are now 'logged
		// in'.
		session := app.Sessions.Load(r)
		err = session.PutInt(w, "currentUserID", currentUserID)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}

	// Redirect the user to the Add Snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the currentUserID from the session data.
	session := app.Sessions.Load(r)
	err := session.Remove(w, "currentUserID")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Redirect the user to the homepage.
	http.Redirect(w, r, "/", 303)
}

func (app *App) SignupAdmin(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.admin.page.html", &HTMLData{
		Form: &forms.SignupAdmin{},
	})
}

func (app *App) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	form := &forms.SignupAdmin{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	if !form.Valid() {
		app.RenderHTML(w, r, "signup.admin.page.html", &HTMLData{Form: form})
		return
	}

	err = app.Database.InsertAdmin(form.Name, form.Email, form.Password)
	if err == models.ErrDuplicateEmail {
		form.Failures["Email"] = "Address is already in use"
		app.RenderHTML(w, r, "signup.admin.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	msg := "Your Admin signup was successful. Please log in using your credentials."
	session := app.Sessions.Load(r)
	err = session.PutString(w, "flash", msg)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}