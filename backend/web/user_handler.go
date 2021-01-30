// The web handler evolving around users, with HTTP-handler functions
// consisting of "GET"- and "POST"-methods. It utilizes session management and
// database access.

package web

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/util"
)

const (
	TokenLength = 43
)

var (
	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	usersRegisterTemplate, usersLoginTemplate, usersProfileTemplate, usersListTemplate, usersForgotPasswordTemplate,
	usersResetPasswordTemplate *template.Template
)

// init gets initialized with the package.
//
// It registers certain types to the session, because by default the session
// can only contain basic data types (int, bool, string, etc.).
//
// All HTML-templates get parsed once to be executed when needed. This is way
// more efficient than parsing the HTML-templates with every request.
func init() {
	gob.Register(x.Token{})

	if _testing { // skip initialization of templates when running tests
		return
	}

	usersRegisterTemplate = template.Must(template.ParseFiles(layout, css, path+"users_register.html"))
	usersLoginTemplate = template.Must(template.ParseFiles(layout, css, path+"users_login.html"))
	usersProfileTemplate = template.Must(template.ParseFiles(layout, css, path+"users_profile.html"))
	usersListTemplate = template.Must(template.ParseFiles(layout, css, path+"users_list.html"))
	usersForgotPasswordTemplate = template.Must(template.ParseFiles(layout, css,
		path+"users_forgot_password.html"))
	usersResetPasswordTemplate = template.Must(template.ParseFiles(layout, css,
		path+"users_reset_password.html"))
}

// UserHandler is the object for handlers to access sessions and database.
type UserHandler struct {
	store    x.Store
	sessions *scs.SessionManager
}

// Register is a GET-method that is accessible to anyone not logged in.
//
// It displays a form, in which values for registering as a new user can be
// entered.
func (h *UserHandler) Register() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			h.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersRegisterTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// RegisterSubmit is a POST-method that is accessible to anyone not logged in
// after Register.
//
// It validates the form from Register and redirects to Register in case of an
// invalid input with corresponding error messages. In case of a valid form, it
// stores the new user in the database and redirects to the home-page.
func (h *UserHandler) RegisterSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := RegisterForm{
			Username:      strings.ToLower(req.FormValue("username")),
			Email:         strings.ToLower(req.FormValue("email")),
			Password:      req.FormValue("password"),
			UsernameTaken: false,
			EmailTaken:    false,
		}

		// Check if username is taken
		_, err := h.store.GetUserByUsername(form.Username)
		if err == nil {
			// If error is nil, a user with that username was found, which
			// means the username is already taken
			form.UsernameTaken = true
		}

		// Check if email is taken
		_, err = h.store.GetUserByEmail(form.Email)
		if err == nil {
			// If error is nil, a user with that email was found, which means
			// the email is already taken
			form.EmailTaken = true
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Encrypt password to hash
		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// New user
		user := x.User{
			Username: form.Username,
			Email:    form.Email,
			Password: string(password),
		}
		// Execute SQL statement to create a user
		if err = h.store.CreateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message
		h.sessions.Put(req.Context(), "flash_success",
			"Willkommen "+form.Username+"! Ihre Registrierung war erfolgreich. Loggen Sie sich bitte ein.\n"+
				"Dazu wurde eine Bestätigungs-Email an "+form.Email+" versandt, um Ihre Email zu verifizieren.")

		// New token
		token := x.Token{
			TokenID: util.GenerateString(TokenLength),
			UserID:  user.UserID,
		}

		// Send email to verify a user's email
		EmailVerificationEmail(user, token.TokenID).Send()

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Login is a GET-method that is accessible to anyone not logged in.
//
// It displays a form in which values for logging in as a returning user can be
// entered.
func (h *UserHandler) Login() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			h.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersLoginTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// LoginSubmit is a POST-method that is accessible to anyone not logged in
// after Login.
//
// It validates the form from Login and redirects to Login in case of an
// invalid input with corresponding error messages. In case of a valid form,
// it stores the user in the session and redirects to the home-page.
func (h *UserHandler) LoginSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := LoginForm{
			UsernameOrEmail:          strings.ToLower(req.FormValue("username")),
			Password:                 req.FormValue("password"),
			IncorrectUsernameOrEmail: false,
			IncorrectPassword:        false,
		}

		// Execute SQL statement to get a user
		user, err := h.store.GetUserByUsername(form.UsernameOrEmail) // check if username is correct
		if err != nil {
			// In case of an error, the username doesn't exist
			// Execute SQL statement to get a user
			user, err = h.store.GetUserByEmail(form.UsernameOrEmail) // check if email is correct

			// In case of an error, the email doesn't exist, which means
			// username and email are both incorrect
			form.IncorrectUsernameOrEmail = err != nil

		}
		if err == nil {
			// If username or email is correct, check if password is correct
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))

			// If error is nil, the password matches the hash, which means it
			// is correct.
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
		h.sessions.Put(req.Context(), "flash_success", "Hallo "+user.Username+"! Sie sind nun eingeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Logout is a GET-method that is accessible to any user.
//
// It removes user from the session and redirects to the home-page.
func (h *UserHandler) Logout() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error", "Sie sind gar nicht eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Remove user ID from session
		h.sessions.Remove(req.Context(), "user_id")

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Sie wurden erfolgreich ausgeloggt.")

		// Redirect to Home
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// Profile is a GET-Method that is accessible to any user.
//
// It displays a user's username and statistics, with the ability to change
// username or password.
func (h *UserHandler) Profile() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		User x.User
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Get user logged in
		userInf := req.Context().Value("user")
		if userInf == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Profil zu betrachten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}
		user := userInf.(x.User)

		// Execute HTML-templates with data
		if err := usersProfileTemplate.Execute(res, data{
			User:        user,
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// List is a GET-method that is accessible to any admin.
//
// It lists all users with the ability to delete a user, to promote a user to
// admin, or to reset a user's password.
func (h *UserHandler) List() http.HandlerFunc {

	// Data to pass to HTML-template
	type data struct {
		SessionData
		CSRF template.HTML

		Users []x.User
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(x.User).Admin {
			// If no user is logged in or user logged in isn't an admin, then
			// redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Sie müssen als Admin eingeloggt sein, um alle Benutzer aufzulisten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to get users
		users, err := h.store.GetUsers()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = usersListTemplate.Execute(res, data{
			Users:       users,
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusFound)
			return
		}
	}
}

// Delete is a POST-method that is accessible to any admin.
//
// It deletes the user and redirects to List.
func (h *UserHandler) Delete() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user ID from URL parameters
		userID, _ := strconv.Atoi(chi.URLParam(req, "userID"))

		// Execute SQL statement to delete a user
		err := h.store.DeleteUser(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of users
		http.Redirect(res, req, "/users", http.StatusFound)
	}
}

// Promote is a POST-method that is accessible to any admin.
//
// It promotes a user to an admin and redirects to List.
func (h *UserHandler) Promote() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user ID from URL parameters
		userID, _ := strconv.Atoi(chi.URLParam(req, "userID"))

		// Execute SQL statement to get user
		user, err := h.store.GetUser(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Make user an admin
		user.Admin = true

		// Execute SQL statement to update user
		if err := h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of users
		http.Redirect(res, req, "/users", http.StatusFound)
	}
}

// VerifyEmail is a POST-method that is accessible to anyone.
//
// It verifies a user's email, which allows the user to later on reset their
// password if necessary and removes the reoccurring flash message when logging
// in.
func (h *UserHandler) VerifyEmail() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve token from URL query
		tokenID := req.URL.Query().Get("token")

		// Execute SQL statement to get token
		token, err := h.store.GetToken(tokenID)
		if err != nil {
			// If token doesn't exist, then redirect to home-page with flash
			// message.
			h.sessions.Put(req.Context(), "flash_error", "Ihr Token zum Bestätigen des Emails ist ungültig.")
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		// Execute SQL statement to get user
		user, err := h.store.GetUser(token.UserID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update user
		user.Verified = true

		// Execute SQL statement to update user
		if err = h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to delete tokens
		if err = h.store.DeleteTokensByUser(token.UserID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Ihre Email wurde erfolgreich bestätigt.")

		// Redirect to home-page
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// ResendVerifyEmail is a POST-method that is accessible to any user.
//
// It resends a email-verification email to the user's email address.
func (h *UserHandler) ResendVerifyEmail() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve user from session
		user := req.Context().Value("user").(x.User)

		// New token
		token := x.Token{
			TokenID: util.GenerateString(TokenLength),
			UserID:  user.UserID,
			Expiry:  time.Now().Add(time.Hour),
		}

		// Execute SQL statement to create a token
		if err := h.store.CreateToken(&token); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Eine Bestätigungs-Email wurde an "+user.Email+"versandt.")

		// Send email to verify a user's email
		EmailVerificationEmail(user, token.TokenID).Send()

		// Redirect to home-page
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// ForgotPassword is a GET-method that is accessible to anyone not logged in.
//
// It displays a form, in which the user can receive an email with a link and
// token to reset the password by specifying his/her email.
func (h *UserHandler) ForgotPassword() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			h.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersForgotPasswordTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ForgotPasswordSubmit is a POST-method that is accessible to any user not
// logged in after ForgotPassword.
//
// It validates the form from ForgotPassword and redirects to ForgotPassword in
// case of an invalid input with corresponding error message. In case of a
// valid form, it stores the new token in the database and redirects to the
// home-page with a flash message.
func (h *UserHandler) ForgotPasswordSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve email from form
		form := ForgotPasswordForm{
			Email: strings.ToLower(req.FormValue("email")),
		}

		// Check if email is valid
		user, err := h.store.GetUserByEmail(form.Email)
		// If error is nil, a user with that email was found.
		form.IncorrectEmail = err != nil
		// If user's email isn't verified, he/she can't reset the password via
		// email
		form.UnverifiedEmail = !user.Verified

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// New token
		token := x.Token{
			TokenID: util.GenerateString(TokenLength),
			UserID:  user.UserID,
			Expiry:  time.Now().Add(time.Hour),
		}

		// Execute SQL statement to create new token
		if err = h.store.CreateToken(&token); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send email to reset a user's password
		PasswordResetEmail(user, token.TokenID).Send()

		// Add flash message to session
		h.sessions.Put(req.Context(), "form",
			"Eine Email zum Zurücksetzen Ihres Passworts wurde an "+form.Email+" versandt.")

		// Redirect to home-page
		http.Redirect(res, req, "/", http.StatusFound)
	}
}

// ResetPassword is a GET-method that is accessible to any user not
// logged in after opening a link with a valid token received via email.
//
// It a verifies the token and displays a form, in which the user can enter a
// new password.
func (h *UserHandler) ResetPassword() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a userInf is logged in
		userInf := req.Context().Value("userInf")
		if userInf != nil {
			// If a user is already logged in, then redirect back with flash
			// message
			h.sessions.Put(req.Context(), "flash_error", "Sie sind bereits eingeloggt.")
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		// Retrieve token from URL query
		tokenID := req.URL.Query().Get("token")

		// Execute SQL statement to get token
		token, err := h.store.GetToken(tokenID)
		if err != nil {
			// If token doesn't exist, then redirect to home-page with flash
			// message.
			h.sessions.Put(req.Context(), "flash_error", "Ihr Token zum Zurücksetzen des Passworts ist ungültig.")
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		// Validate expiration time of token
		if token.Expiry.After(time.Now()) {
			// If the token has expired (after 1 hour), then redirect back with
			// flash message
			h.sessions.Put(req.Context(), "flash_error", "Der Token ist abgelaufen. Sie haben jeweils 1 Stunde Zeit, "+
				"um Ihr Passwort zurückzusetzen.")
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		// Add token to session
		h.sessions.Put(req.Context(), "token", token)

		// Execute HTML-templates with data
		if err = usersResetPasswordTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ResetPasswordSubmit is a POST-method that is accessible to anyone not logged
// in.
//
// It changes a user's password, and deletes all of the user's tokens.
func (h *UserHandler) ResetPasswordSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve new password from form
		form := ResetPasswordForm{
			Password: req.FormValue("password"),
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve token from session
		token := h.sessions.Get(req.Context(), "token").(x.Token)

		// Execute SQL statement to get user
		user, err := h.store.GetUser(token.UserID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encrypt password to hash
		password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update user's password
		user.Password = string(password)

		// Execute SQL statement to update user
		if err = h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to delete tokens
		if err = h.store.DeleteTokensByUser(user.UserID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Ihr Passwort wurde erfolgreich geändert. "+
			"Bitte loggen Sie sich ein.")

		// Redirect to login
		http.Redirect(res, req, "/users/login", http.StatusNotFound)
	}
}
