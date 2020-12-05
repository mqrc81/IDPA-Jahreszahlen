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

/*
 * Register is a GET method with form to register new user
 */
func (h *UserHandler) Register() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(usersRegisterHTML))

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

		// Validate form
		if _, err := h.store.UserByUsername(form.Username); err == nil { // If error is nil, user was found
			form.UsernameTaken = true // If user was found, username is already taken
		}
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
			UserID:   0,
			Username: form.Username,
			Password: string(password),
			Admin:    false,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message
		h.sessions.Put(req.Context(), "flash", "Registrierung war erfolgreich. Bitte loggen Sie sich ein.")

		http.Redirect(res, req, "/users/login", http.StatusFound)
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
	tmpl := template.Must(template.New("").Parse(usersLoginHTML))

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
			Username:             req.FormValue("username"),
			Password:             req.FormValue("password"),
			IncorrectCredentials: false,
		}

		user, err := h.store.UserByUsername(form.Username)
		if err != nil {
			form.IncorrectCredentials = true
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
			form.IncorrectCredentials = err != nil
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
		h.sessions.Put(req.Context(), "flash", "Login war erfolgreich.")

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
		h.sessions.Put(req.Context(), "flash", "Logout war erfolgreich.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}
