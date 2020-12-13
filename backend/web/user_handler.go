package web

/*
 * user_handler.go contains HTTP-handler functions for users
 */

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * UserHandler handles sessions, CSRF-protection and database access for users
 */
type UserHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// TODO
// Profile()
// EditUsername()
// EditUsernameStore()

/*
 * Register is a GET method with form to register new user
 */
func (h *UserHandler) Register() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_register.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * RegisterSubmit is a POST method that stores user created
 */
func (h *UserHandler) RegisterSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form (Register)
		form := RegisterForm{
			Username:      req.FormValue("username"),
			Password:      req.FormValue("password"),
			UsernameTaken: false,
		}

		// Check if username is taken
		user, err := h.store.UserByUsername(form.Username)
		if err == nil { // If error is nil, user was found
			form.UsernameTaken = true // If user was found, username is already taken
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Hash password
		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement
		if err := h.store.CreateUser(&backend.User{
			Username: form.Username,
			Password: string(password),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store user ID in session (= login)
		h.sessions.Put(req.Context(), "user_id", user.UserID)

		// Add flash message
		h.sessions.Put(req.Context(), "flash", "Willkommen "+form.Username+"! Ihre Registrierung war erfolgreich. "+
			"Sie sind nun eingeloggt.")

		http.Redirect(res, req, "/", http.StatusFound)
	}
}

/*
 * Login is a GET method with form to login
 */
func (h *UserHandler) Login() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_login.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * LoginSubmit is a POST method that logs in user
 */
func (h *UserHandler) LoginSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		form := LoginForm{
			Username:          req.FormValue("username"),
			Password:          req.FormValue("password"),
			IncorrectUsername: false,
			IncorrectPassword: false,
		}

		// Execute SQL statement
		user, err := h.store.UserByUsername(form.Username)
		if err != nil {
			// In case of an error, the username doesn't exist
			form.IncorrectUsername = true
		} else {
			// Else, check if username and password match
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
			form.IncorrectPassword = err != nil
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Store user ID in session
		h.sessions.Put(req.Context(), "user_id", user.UserID)

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash", "Hallo "+form.Username+"! Sie sind nun eingeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

/*
 * Logout is a POST method that logs out user
 */
func (h *UserHandler) Logout() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Remove user ID from session
		h.sessions.Remove(req.Context(), "user_id")

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash", "Sie wurden erfolgreich ausgeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

/*
 * EditPassword is a GET method that edits a user's password
 */
func (h *UserHandler) EditPassword() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_edit_password.html"))
	return func(res http.ResponseWriter, req *http.Request) {

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * EditPasswordStore is a POST method that stores password edited
 */
func (h *UserHandler) EditPasswordStore() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		oldPassword := req.FormValue("old_password")
		form := PasswordForm{
			Password1:            req.FormValue("password1"),
			Password2:            req.FormValue("password2"),
			IncorrectOldPassword: false,
		}

		// Retrieve user from session
		var userID int
		userIDinf := h.sessions.Get(req.Context(), "user_id")
		if userIDinf != nil {
			userID = userIDinf.(int)
		}
		user, err := h.store.User(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Compare user's password with "old password" from form
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
			form.IncorrectOldPassword = true
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement
		if err := h.store.UpdateUser(&backend.User{
			UserID:   userID,
			Username: user.Username,
			Password: form.Password1,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
