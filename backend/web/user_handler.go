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
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(usersRegisterHTML))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(res, nil); err != nil {
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
		// Hash password
		password, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement
		if err := h.store.CreateUser(&backend.User{
			UserID:   0,
			Username: req.FormValue("username"),
			Password: string(password),
			Admin:    false,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(res, req, "/users/login", http.StatusFound)
	}
}

/*
 * Login is a GET method with form to login
 */
func (h *UserHandler) Login() http.HandlerFunc {
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(usersLoginHTML))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(res, nil); err != nil {
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
		// Execute SQL statement
		user, err := h.store.UserByUsername(req.FormValue("username"))
		if err != nil {
			// TODO username incorrect
		} else {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))); err != nil {
				// TODO password incorrect
			}
		}
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

/*
 * Logout is a POST method that logs out user
 */
func (h *UserHandler) Logout() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// TODO log out

		http.Redirect(res, req, "/", http.StatusFound)
	}
}
