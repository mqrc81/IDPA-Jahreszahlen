// A branch of the user handler (for a better overview), which contains HTTP-
// handlers that deal with editing a user's username, email and password.

package web

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	usersEditUsernameTemplate, usersEditEmailTemplate, usersEditPasswordTemplate *template.Template
)

// init gets initialized with the package.
//
// All HTML-templates get parsed once to be executed when needed. This is way
// more efficient than parsing the HTML-templates with every request.
func init() {
	if _testing { // skip initialization of templates when running tests
		return
	}

	usersEditUsernameTemplate = template.Must(template.ParseFiles(layout, css, path+"users_edit_username.html"))
	usersEditEmailTemplate = template.Must(template.ParseFiles(layout, css, path+"users_edit_email.html"))
	usersEditPasswordTemplate = template.Must(template.ParseFiles(layout, css, path+"users_edit_password.html"))
}

// EditUsername is a GET-method that is accessible to any user.
//
// It displays a form in which values for modifying the current username can be
// entered.
func (h *UserHandler) EditUsername() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Benutzernamen zu ändern.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersEditUsernameTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditUsernameSubmit is a POST-method that is accessible to any user after
// EditUsername.
//
// It validates the form from EditUsername and redirects to EditUsername in
// case of an invalid input with corresponding error messages. In case of a
// valid form, it stores the user in the database and redirects to Profile.
func (h *UserHandler) EditUsernameSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := EditUsernameForm{
			NewUsername: strings.ToLower(req.FormValue("username")),
			Password:    req.FormValue("password"),
		}

		// Check if username is taken
		_, err := h.store.GetUserByUsername(form.NewUsername)
		// If error is nil, a user with that username was found, which means
		// the username is already taken.
		form.UsernameTaken = err == nil

		// Retrieve user from session
		user := req.Context().Value("user").(x.User)

		// Check if password is correct
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		// If error is nil, the password matches the hash, which means it is
		// correct
		form.IncorrectPassword = err != nil

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Update user's username
		user.Username = form.NewUsername

		// Execute SQL statement to update user
		if err = h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Ihr Benutzername wurde erfolgreich geändert.")

		// Redirect to profile
		http.Redirect(res, req, "/users/profile", http.StatusFound)
	}
}

// EditEmail is a GET-method that is accessible to any user.
//
// It displays a form in which values for modifying the current email can be
// entered.
func (h *UserHandler) EditEmail() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihre Email zu ändern.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersEditEmailTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditEmailSubmit is a POST-method that is accessible to any user after
// EditEmail.
//
// It validates the form from EditEmail and redirects to EditEmail in
// case of an invalid input with corresponding error messages. In case of a
// valid form, it stores the user in the database and redirects to Profile.
func (h *UserHandler) EditEmailSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := EditEmailForm{
			NewEmail: strings.ToLower(req.FormValue("email")),
			Password: req.FormValue("password"),
		}

		// Check if email is taken
		_, err := h.store.GetUserByEmail(form.NewEmail)
		// If error is nil, a user with that username was found, which means
		// the email is already taken.
		form.EmailTaken = err == nil

		// Retrieve user from session
		user := req.Context().Value("user").(x.User)

		// Check if password is correct
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		// If error is nil, the password matches the hash, which means it is
		// correct
		form.IncorrectPassword = err != nil

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Update user's email
		user.Email = form.NewEmail

		// Execute SQL statement to update user
		if err = h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Ihre Email wurde erfolgreich geändert.")

		// Redirect to profile
		http.Redirect(res, req, "/users/profile", http.StatusFound)
	}
}

// EditPassword is a GET-method that is accessible to any user.
//
// It displays a form in which values for modifying the current password can
// be entered.
func (h *UserHandler) EditPassword() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Loggen Sie sich zuerst ein, um Ihr Passwort zu ändern.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := usersEditPasswordTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditPasswordSubmit is a POST-method that is accessible to any user after
// EditPassword.
//
// It validates the form from EditPassword and redirects to EditPassword in
// case of an invalid input with corresponding error messages. In case of a
// valid form, it stores the user in the database and redirects to Profile.
func (h *UserHandler) EditPasswordSubmit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := EditPasswordForm{
			NewPassword: req.FormValue("new_password"),
			Password:    req.FormValue("password"),
		}

		// Retrieve user from session
		user := req.Context().Value("user").(x.User)

		// Compare user's password with "old password" from form
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
			form.IncorrectPassword = true
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Encrypt password to hash
		password, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update user's password
		user.Password = string(password)

		// Execute SQL statement to update a user
		if err := h.store.UpdateUser(&user); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Ihr Passwort wurde erfolgreich geändert.")

		// Redirect to user's profile
		http.Redirect(res, req, "/users/profile", http.StatusFound)
	}
}
