package web

// forms.go
// Contains form validation for all necessary HTTP-handlers.

import (
	"encoding/gob"
	"log"
	"regexp"
	"time"
)

// init
// Gets initialized with the package. Registers certain types to the session,
// because by default the session can only contain basic data types (int, bool,
// string, etc.).
func init() {
	gob.Register(TopicForm{})
	gob.Register(EventForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
	gob.Register(PasswordForm{})
	gob.Register(FormErrors{})
}

// FormErrors
// A map that holds the error messages. The key string contains the name of the
// form input that is invalid (e.g. "Title"), the value string contains the
// error message.
type FormErrors map[string]string

// TopicForm
// Holds values of the form input when creating or editing a
// topic.
type TopicForm struct {
	Title       string
	StartYear   int
	EndYear     int
	Description string

	Errors FormErrors
}

// Validate
// Validates the form input when creating or editing a topic.
func (f *TopicForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate title
	if f.Title == "" {
		f.Errors["Title"] = "Titel darf nicht leer sein."
	} else if len(f.Title) > 50 {
		f.Errors["Title"] = "Titel darf 50 Zeichen nicht überschreiten."
	}

	// Validate start- and end-year
	current := time.Now().Year()
	if f.StartYear <= 0 {
		f.Errors["Year"] = "Start-Jahr muss positiv sein."
	} else if f.EndYear <= 0 {
		f.Errors["Year"] = "End-Jahr muss positiv sein."
	} else if f.StartYear > current {
		f.Errors["Year"] = "Start-Jahr darf nicht in der Zukunft sein."
	} else if f.EndYear > current {
		f.Errors["Year"] = "End-Jahr darf nicht in der Zukunft sein."
	} else if f.EndYear < f.StartYear {
		f.Errors["Year"] = "Da wurden wohl Start- und End-Jahr vertauscht."
	}

	// Validate description
	if len(f.Description) > 500 {
		f.Errors["Description"] = "Beschreibung darf nicht leer sein."
	}

	return len(f.Errors) == 0
}

// EventForm
// Holds values of the form input when creating or editing an
// event.
type EventForm struct {
	Title string
	Year  int

	Errors FormErrors
}

// Validate
// Validates the form input when creating or editing an event.
func (f *EventForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate title
	if f.Title == "" {
		f.Errors["Title"] = "Titel darf nicht leer sein."
	} else if len(f.Title) > 110 {
		f.Errors["Title"] = "Titel darf 110 Zeichen nicht überschreiten."
	}

	// Validate year
	if f.Year == 0 {
		f.Errors["Year"] = "Jahr darf nicht leer sein."
	} else if f.Year <= 0 {
		f.Errors["Year"] = "Jahr muss positiv sein."
	} else if f.Year > time.Now().Year() {
		f.Errors["Year"] = "Wird hier die Zukunft vorausgesagt?"
	}

	return len(f.Errors) == 0
}

// RegisterForm
// Holds values of the form input when registering.
type RegisterForm struct {
	Username      string
	Password      string
	UsernameTaken bool

	Errors FormErrors
}

// Validate
// Validates the form input when registering.
func (f *RegisterForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate new password
	if f.UsernameTaken {
		f.Errors["Username"] = "Benutzername ist bereits vergeben."
	} else {
		f.Errors = validateUsername(f.Username, f.Errors)
	}

	// Validate password
	f.Errors = validatePassword(f.Password, f.Errors, "Password")

	return len(f.Errors) == 0
}

// LoginForm
// Holds values of the form input when logging in.
type LoginForm struct {
	Username          string
	Password          string
	IncorrectUsername bool
	IncorrectPassword bool

	Errors FormErrors
}

// Validate
// Validates the form input when logging in.
func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate username
	if f.Username == "" {
		f.Errors["Username"] = "Bitte Benutzernamen eingeben."
	} else if f.IncorrectUsername {
		f.Errors["Username"] = "Ungültiger Benutzername."
	}

	// Validate password
	if f.Password == "" {
		f.Errors["Password"] = "Bitte Passwort eingeben."
	} else if f.IncorrectPassword {
		f.Errors["Password"] = "Ungültiges Passwort."
	}

	return len(f.Errors) == 0
}

// UsernameForm
// Holds values of the form input when editing a username.
type UsernameForm struct {
	NewUsername       string
	Password          string
	UsernameTaken     bool
	IncorrectPassword bool

	Errors FormErrors
}

// Validate
// Validates the form input when editing a password.
func (f *UsernameForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate new username
	if f.UsernameTaken {
		f.Errors["Username"] = "Benutzername ist bereits vergeben."
	} else {
		f.Errors = validateUsername(f.NewUsername, f.Errors)
	}

	// Validate password
	if f.Password == "" {
		f.Errors["OldPassword"] = "Geben Sie Ihr Passwort ein."
	} else if f.IncorrectPassword {
		f.Errors["OldPassword"] = "Passwort ist inkorrekt."
	}

	return len(f.Errors) == 0
}

// PasswordForm
// Holds values of the form input when editing a password.
type PasswordForm struct {
	NewPassword          string
	OldPassword          string
	IncorrectOldPassword bool

	Errors FormErrors
}

// Validate
// Validates the form input when editing a password.
func (f *PasswordForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate old password
	if f.OldPassword == "" {
		f.Errors["OldPassword"] = "Geben Sie Ihr altes Passwort ein."
	} else if f.IncorrectOldPassword {
		f.Errors["OldPassword"] = "Altes Passwort ist inkorrekt."
	}

	// Validate new password
	f.Errors = validatePassword(f.NewPassword, f.Errors, "NewPassword")

	return len(f.Errors) == 0
}

// validateUsername
// Validates a username.
func validateUsername(username string, errors FormErrors) FormErrors {
	if len(username) < 3 {
		errors["Username"] = "Benutzername muss mindestens 3 Zeichen lang sein."
	} else if len(username) > 20 {
		errors["Username"] = "Benutzername darf höchstens 20 Zeichen lang sein."
	} else if !regex(username, "^[a-zA-Z0-9._]*$") {
		errors["Username"] = "Benutzername darf nur Buchstaben, Zahlen, '.' und '_' enthalten."
	} else if !regex(username, "\\D") {
		errors["Username"] = "Benutzername muss mindestens 1 Buchstaben enthalten."
	} else if regex(username, "^[._]") {
		errors["Username"] = "Benutzername darf nicht mit '.' oder '_' beginnen."
	} else if regex(username, "[._]$") {
		errors["Username"] = "Benutzername darf nicht mit '.' oder '_' enden."
	} else if regex(username, "[_.]{2}") {
		errors["Username"] = "Benutzername darf '.' und '_' nicht aufeinanderfolgend haben."
	}

	return errors
}

// validatePassword
// Validates a password.
func validatePassword(password string, errors FormErrors, errorName string) FormErrors {
	if len(password) < 8 {
		errors[errorName] = "Passwort muss mindestens 8 Zeichen lang sein."
	} else if !regex(password, "[!@#$%^&*]") {
		errors[errorName] = "Passwort muss ein Sonderzeichen enthalten (!@#$%^&*)."
	} else if !regex(password, "[a-z]") {
		errors[errorName] = "Passwort muss mindestens ein Kleinbuchstaben enthalten."
	} else if !regex(password, "[A-Z]") {
		errors[errorName] = "Passwort muss mindestens ein Grossbuchstaben enthalten."
	} else if !regex(password, "\\d") {
		errors[errorName] = "Passwort muss mindestens eine Zahl enthalten."
	}

	return errors
}

// regex
// Checks if a certain regular expression matches a certain string.
func regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(err)
	}
	return match
}
