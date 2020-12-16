package web

// user_handler.go
// Contains all HTTP-handlers for pages evolving around users.

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// UserHandler
// Object for handlers to access sessions and database.
type UserHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

//TODO
// Profile()

// Register
// A GET-method. It renders a form, in which values for registering can be
// entered.
func (handler *UserHandler) Register() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_register.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterSubmit
// A POST-method. It validates the form from Register and redirects to Register
// in case of an invalid input with corresponding error messages. In case of a
// valid form, it stores the new user in the database and redirects to the home-
// page.
func (handler *UserHandler) RegisterSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form (Register)
		form := RegisterForm{
			Username:      strings.ToLower(req.FormValue("username")),
			Password:      req.FormValue("password"),
			UsernameTaken: false,
		}

		// Check if username is taken
		_, err := handler.store.UserByUsername(form.Username)
		if err == nil {
			// If error is nil, a user with that username was found, which
			// means the username is already taken
			form.UsernameTaken = true
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Hash password
		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to create a user
		if err := handler.store.CreateUser(&backend.User{
			Username: form.Username,
			Password: string(password),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message
		handler.sessions.Put(req.Context(), "flash", "Willkommen "+form.Username+"! Ihre Registrierung war erfolgreich. "+
			"Loggen Sie sich bitte ein.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Login
// A GET-method. It renders a form in which values for logging in can be
// entered.
func (handler *UserHandler) Login() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_login.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// LoginSubmit
// A POST-method. It validates the form from Login and redirects to Login
// in case of an invalid input with corresponding error messages. In case of a
// valid form, it stores the user in the session and redirects to the home-page.
func (handler *UserHandler) LoginSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		form := LoginForm{
			Username:          strings.ToLower(req.FormValue("username")),
			Password:          req.FormValue("password"),
			IncorrectUsername: false,
			IncorrectPassword: false,
		}

		// Execute SQL statement to get a user
		user, err := handler.store.UserByUsername(form.Username)
		if err != nil {
			// In case of an error, the username doesn't exist
			form.IncorrectUsername = true
		} else {
			// Else, check if password is correct
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
			// If error is nil, the password matches the hash, which means it
			// is correct.
			form.IncorrectPassword = err != nil
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// CreateStore user ID in session
		handler.sessions.Put(req.Context(), "user_id", user.UserID)

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash", "Hallo "+form.Username+"! Sie sind nun eingeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Logout
// A GET-method that any user can call. It removes user from the session.
func (handler *UserHandler) Logout() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Remove user ID from session
		handler.sessions.Remove(req.Context(), "user_id")

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash", "Sie wurden erfolgreich ausgeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// EditUsername
// A GET-method that any user can call. It renders a form in which values for
// updating the current username can be entered.
func (handler *UserHandler) EditUsername() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_edit_username.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditUsernameSubmit
// A POST-method. It validates the form from EditUsername and redirects to
// EditUsername in case of an invalid input with corresponding error messages.
// In case of a valid form, it stores the user in the database and redirects to
// the user's profile.
func (handler *UserHandler) EditUsernameSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		form := UsernameForm{
			NewUsername:       strings.ToLower(req.FormValue("username")),
			Password:          req.FormValue("password"),
			UsernameTaken:     false,
			IncorrectPassword: false,
		}

		// Check if username is taken
		_, err := handler.store.UserByUsername(form.NewUsername)
		// If error is nil, a user with that username was found, which means
		// the username is already taken.
		form.UsernameTaken = err == nil

		// Retrieve user from session
		user := req.Context().Value("user").(backend.User)

		// Check if password is correct
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		// If error is nil, the password matches the hash, which means it is
		// correct
		form.IncorrectPassword = err != nil

		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// CreateStore user ID in session
		handler.sessions.Put(req.Context(), "user_id", user.UserID)

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash", "Ihr Benutzername wurde erfolgreich geändert.")

		// Redirect to Home
		http.Redirect(res, req, "/profile", http.StatusFound)
	}
}

// EditPassword
// A GET-method that any user can call. It renders a form in which values for
// updating the current password can be entered.
func (handler *UserHandler) EditPassword() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/users_edit_password.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditPasswordSubmit
// A POST-method. It validates the form from EditPassword and redirects to
// EditPassword in case of an invalid input with corresponding error messages.
// In case of a valid form, it stores the user in the database and redirects to
// the user's profile.
func (handler *UserHandler) EditPasswordSubmit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		form := PasswordForm{
			NewPassword:          req.FormValue("new_password"),
			OldPassword:          req.FormValue("old_password"),
			IncorrectOldPassword: false,
		}

		// Retrieve user from session
		user := req.Context().Value("user").(backend.User)

		// Compare user's password with "old password" from form
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.OldPassword)); err != nil {
			form.IncorrectOldPassword = true
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Hash password
		password, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to update a user
		if err := handler.store.UpdateUser(&backend.User{
			UserID:   user.UserID,
			Username: user.Username,
			Password: string(password),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to user's profile
		http.Redirect(res, req, "/users/profile", http.StatusFound)
	}
}
