package web

/*
 * Contains all HTTP-handler functions for pages evolving around users.
 */

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// UserHandler is the object for handlers to access sessions and database.
type UserHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// Register is a GET-method. It renders a form, in which values for registering
// can be entered.
func (handler *UserHandler) Register() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_register.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			handler.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterSubmit is a POST-method. It validates the form from Register and
// redirects to Register in case of an invalid input with corresponding error
// messages. In case of a valid form, it stores the new user in the database
// and redirects to the home-page.
func (handler *UserHandler) RegisterSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form (Register)
		form := RegisterForm{
			Username:      strings.ToLower(req.FormValue("username")),
			Password:      req.FormValue("password"),
			UsernameTaken: false,
		}

		// Check if username is taken
		_, err := handler.store.GetUserByUsername(form.Username)
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
		handler.sessions.Put(req.Context(), "flash_success",
			"Willkommen "+form.Username+"! Ihre Registrierung war erfolgreich. Loggen Sie sich bitte ein.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Login is a GET-method. It renders a form in which values for logging in can
// be entered.
func (handler *UserHandler) Login() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_login.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			handler.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// LoginSubmit is a POST-method. It validates the form from Login and redirects
// to Login in case of an invalid input with corresponding error messages. In
// case of a valid form, it stores the user in the session and redirects to the
// home-page.
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
		user, err := handler.store.GetUserByUsername(form.Username)
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

		// Store user ID in session
		handler.sessions.Put(req.Context(), "user_id", user.UserID)

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash_success", "Hallo "+form.Username+"! Sie sind nun eingeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Logout is a GET-method that any user can call. It removes user from the
// session.
func (handler *UserHandler) Logout() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Sie sind gar nicht eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Remove user ID from session
		handler.sessions.Remove(req.Context(), "user_id")

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash_success", "Sie wurden erfolgreich ausgeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Profile is a GET-Method that displays a user's username and statistics, with
// the options to change username or password.
func (handler *UserHandler) Profile() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		User backend.User

		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_profile.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Get user logged in
		userInf := req.Context().Value("user")
		if userInf == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Profil zu betrachten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}
		user := userInf.(backend.User)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			User:        user,
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// List is a GET-method that any admin can call. It lists all users with the
// options to delete or promote a user or reset a user's password.
func (handler *UserHandler) List() http.HandlerFunc {

	// Data to pass to HTML-template
	type data struct {
		Users []backend.User

		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_list.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or user logged in isn't an admin, then
			// redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Sie müssen als Admin eingeloggt sein, um alle Benutzer aufzulisten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to get users
		users, err := handler.store.GetUsers()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			Users:       users,
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusFound)
			return
		}
	}
}

// EditUsername is a GET-method that any user can call. It renders a form in
// which values for updating the current username can be entered.
func (handler *UserHandler) EditUsername() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_edit_username.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Benutzernamen zu ändern.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditUsernameSubmit is a POST-method. It validates the form from EditUsername
// and redirects to EditUsername in case of an invalid input with corresponding
// error messages. In case of a valid form, it stores the user in the database
// and redirects to the Profile.
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
		_, err := handler.store.GetUserByUsername(form.NewUsername)
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
		handler.sessions.Put(req.Context(), "flash_success", "Ihr Benutzername wurde erfolgreich geändert.")

		// Redirect to Home
		http.Redirect(res, req, "/profile", http.StatusFound)
	}
}

// EditPassword is a GET-method that any user can call. It renders a form in
// which values for updating the current password can be entered.
func (handler *UserHandler) EditPassword() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/users_edit_password.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Passwort zu ändern.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditPasswordSubmit is a POST-method. It validates the form from EditPassword
// and redirects to EditPassword in case of an invalid input with corresponding
// error messages. In case of a valid form, it stores the user in the database
// and redirects to the Profile.
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

// Delete is a POST-method that any admin can call. It deletes the user and
// redirects to List.
func (handler *UserHandler) Delete() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user ID from URL parameters
		userID, _ := strconv.Atoi(chi.URLParam(req, "userID"))

		// Execute SQL statement to delete a user
		err := handler.store.DeleteUser(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Promote is a POST-method that any admin can call. It promotes a user to an
// admin.
func (handler *UserHandler) Promote() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user ID from URL parameters
		userID, _ := strconv.Atoi(chi.URLParam(req, "userID"))

		// Execute SQL statement to get user
		user, err := handler.store.GetUser(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Make user an admin
		user.Admin = true

		// Execute SQL statement to update user
		if err := handler.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of users
		http.Redirect(res, req, "/users", http.StatusFound)
	}
}

// ResetPassword is a POST-method that any admin can call. It resets a user's
// password.
func (handler *UserHandler) ResetPassword() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user ID from URL parameters
		userID, _ := strconv.Atoi(chi.URLParam(req, "userID"))

		// Execute SQL statement to get user
		user, err := handler.store.GetUser(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Hash password
		password, err := bcrypt.GenerateFromPassword([]byte(DefaultPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Change user's password
		user.Password = string(password)

		// Execute SQL statement to update user
		if err := handler.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of users
		http.Redirect(res, req, "/users", http.StatusFound)
	}
}
