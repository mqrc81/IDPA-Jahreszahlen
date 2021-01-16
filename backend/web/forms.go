// The collection of form validation for various HTTP-handler functions. Once a
// form gets submitted, it gets validated here. The form gets returned,
// including potential error messages to be displayed in the form, after
// being redirected back to the form.

package web

import (
	"encoding/gob"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

// init gets initialized with the package.
//
// It registers certain types to the
// session, because by default the session can only contain basic data types
// (int, bool, string, etc.).
func init() {
	gob.Register(TopicForm{})
	gob.Register(EventForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
	gob.Register(UsernameForm{})
	gob.Register(EmailForm{})
	gob.Register(PasswordForm{})
	gob.Register(ForgotPasswordForm{})
	gob.Register(FormErrors{})
}

// FormErrors is a map that holds the error messages. The key string contains
// the name of the form input that is invalid (e.g. "Name"), the value string
// contains the error message.
type FormErrors map[string]string

// TopicForm holds values of the form input when creating or editing a topic.
type TopicForm struct {
	Name        string
	StartYear   int
	EndYear     int
	Description string

	Errors FormErrors
}

// Validate validates the form input when creating or editing a topic.
func (form *TopicForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate name
	if form.Name == "" {
		form.Errors["Name"] = "Titel darf nicht leer sein."
	} else if len(form.Name) > 50 {
		form.Errors["Name"] = "Titel darf 50 Zeichen nicht überschreiten."
	}

	// Validate start- and end-year
	current := time.Now().Year()
	if form.StartYear <= 0 {
		form.Errors["Year"] = "Start-Jahr muss positiv sein."
	} else if form.EndYear <= 0 {
		form.Errors["Year"] = "End-Jahr muss positiv sein."
	} else if form.StartYear > current {
		form.Errors["Year"] = "Start-Jahr darf nicht in der Zukunft sein."
	} else if form.EndYear > current {
		form.Errors["Year"] = "End-Jahr darf nicht in der Zukunft sein."
	} else if form.EndYear < form.StartYear {
		form.Errors["Year"] = "Da wurden wohl Start- und End-Jahr vertauscht."
	}

	// Validate description
	if len(form.Description) > 500 {
		form.Errors["Description"] = "Beschreibung darf 500 Buchstaben nicht überschreiten."
	}

	return len(form.Errors) == 0
}

// EventForm holds values of the form input when creating or editing an event.
type EventForm struct {
	Name       string
	Year       int
	Date       time.Time
	YearOrDate string

	Errors FormErrors
}

// Validate validates the form input when creating or editing an event.
func (form *EventForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate name
	if form.Name == "" {
		form.Errors["Name"] = "Titel darf nicht leer sein."
	} else if len(form.Name) > 110 {
		form.Errors["Name"] = "Titel darf 110 Zeichen nicht überschreiten."
	}

	// Validate date or year
	year, err := strconv.Atoi(form.YearOrDate) // convert to int; if error occurs, admin entered a date, not a year
	if err == nil {                            // admin entered a year
		form.Year = year
		form.Date, _ = time.Parse("2006", form.YearOrDate) // date = year + default values (e.g. 1969-01-01 00:00:00)

		// Validate year
		if form.Year == 0 {
			form.Errors["Year"] = "Jahr darf nicht leer sein."
		} else if form.Year <= 0 {
			form.Errors["Year"] = "Jahr muss positiv sein."
		} else if form.Year > time.Now().Year() {
			form.Errors["Year"] = "Wird hier die Zukunft vorausgesagt?"
		}
	} else { // admin entered (day, ) month & year (e.g. 08.1969 or 13.19.69)

		date, err := time.Parse("02.01.2006", form.YearOrDate) // check if admin entered valid date as 'dd.mm.yyyy'
		if err != nil {
			date, err = time.Parse("01.2006", form.YearOrDate) // check if admin entered valid date as 'mm.yy'
			if err != nil {
				now := time.Now()
				form.Errors["Year"] = fmt.Sprintf("Ungültiges Format. Erlaubte Formate: '%v', '%s', '%s'",
					now.Year(), now.Format("01.2006"), now.Format("02.01.2006"))
			}
		}

		if form.Errors["Year"] == "" { // admin entered valid date
			if date.Unix() > time.Now().Unix() {
				form.Errors["Year"] = "Wird hier die Zukunft vorausgesagt?"
			}
			form.Date = date
			form.Year = date.Year()
		}
	}

	return len(form.Errors) == 0
}

// RegisterForm holds values of the form input when registering.
type RegisterForm struct {
	Username      string
	Email         string
	Password      string
	UsernameTaken bool
	EmailTaken    bool

	Errors FormErrors
}

// Validate validates the form input when registering.
func (form *RegisterForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate username
	if form.UsernameTaken {
		form.Errors["Username"] = "Benutzername ist bereits vergeben."
	} else {
		form.Errors.validateUsername(form.Username)
	}

	// Validate email
	if form.EmailTaken {
		form.Errors["Email"] = "Email ist bereits vergeben."
	} else {
		form.Errors.validateEmail(form.Email)
	}

	// Validate password
	form.Errors.validatePassword(form.Password, "Password")

	return len(form.Errors) == 0
}

// LoginForm holds values of the form input when logging in.
type LoginForm struct {
	UsernameOrEmail          string
	Password                 string
	IncorrectUsernameOrEmail bool
	IncorrectPassword        bool

	Errors FormErrors
}

// Validate validates the form input when logging in.
func (form *LoginForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate username or email
	if form.UsernameOrEmail == "" {
		form.Errors["UsernameOrEmail"] = "Bitte Benutzernamen order Email angeben."
	} else if form.IncorrectUsernameOrEmail {
		form.Errors["UsernameOrEmail"] = "Ungültiger Benutzername oder Email."
	}

	// Validate password
	if form.Password == "" {
		form.Errors["Password"] = "Bitte Passwort angeben."
	} else if form.IncorrectPassword {
		form.Errors["Password"] = "Ungültiges Passwort."
	}

	return len(form.Errors) == 0
}

// UsernameForm holds values of the form input when editing a username.
type UsernameForm struct {
	NewUsername       string
	Password          string
	UsernameTaken     bool
	IncorrectPassword bool

	Errors FormErrors
}

// Validate validates the form input when editing a password.
func (form *UsernameForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate new username
	if form.UsernameTaken {
		form.Errors["Username"] = "Benutzername ist bereits vergeben."
	} else {
		form.Errors.validateUsername(form.NewUsername)
	}

	// Validate password
	if form.Password == "" {
		form.Errors["OldPassword"] = "Geben Sie Ihr Passwort ein."
	} else if form.IncorrectPassword {
		form.Errors["OldPassword"] = "Passwort ist inkorrekt."
	}

	return len(form.Errors) == 0
}

// EmailForm holds values of the form input when editing an email.
type EmailForm struct {
	NewEmail          string
	Password          string
	EmailTaken        bool
	IncorrectPassword bool

	Errors FormErrors
}

// Validate validates the form input when editing a password.
func (form *EmailForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate new username
	if form.EmailTaken {
		form.Errors["Email"] = "Email ist bereits vergeben."
	} else {
		form.Errors.validateEmail(form.NewEmail)
	}

	// Validate password
	if form.Password == "" {
		form.Errors["OldPassword"] = "Geben Sie Ihr Passwort ein."
	} else if form.IncorrectPassword {
		form.Errors["OldPassword"] = "Passwort ist inkorrekt."
	}

	return len(form.Errors) == 0
}

// PasswordForm holds values of the form input when editing a password.
type PasswordForm struct {
	NewPassword          string
	OldPassword          string
	IncorrectOldPassword bool

	Errors FormErrors
}

// Validate validates the form input when editing a password.
func (form *PasswordForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate old password
	if form.OldPassword == "" {
		form.Errors["OldPassword"] = "Geben Sie Ihr altes Passwort ein."
	} else if form.IncorrectOldPassword {
		form.Errors["OldPassword"] = "Altes Passwort ist inkorrekt."
	}

	// Validate new password
	form.Errors.validatePassword(form.NewPassword, "NewPassword")

	return len(form.Errors) == 0
}

type ForgotPasswordForm struct {
	Email           string
	IncorrectEmail  bool
	UnverifiedEmail bool

	Errors FormErrors
}

func (form *ForgotPasswordForm) Validate() bool {
	form.Errors = FormErrors{}

	// Validate email
	if form.Email == "" {
		form.Errors["Email"] = "Bitte Email angeben."
	} else if form.IncorrectEmail {
		form.Errors["Email"] = "Ungültige Email."
	} else if form.UnverifiedEmail {
		form.Errors["Email"] = "Ihre Email wurde nie bestätigt. Sie können derzeit das Passwort nicht zurücksetzen."
	}

	return len(form.Errors) == 0
}

// validateUsername validates a username.
func (errors *FormErrors) validateUsername(username string) {
	if username == "" {
		(*errors)["Username"] = "Bitte Benutzernamen angeben."
	} else if len(username) < 3 {
		(*errors)["Username"] = "Benutzername muss mindestens 3 Zeichen lang sein."
	} else if len(username) > 20 {
		(*errors)["Username"] = "Benutzername darf höchstens 20 Zeichen lang sein."
	} else if !regex(username, "^[a-zA-Z0-9._]*$") {
		(*errors)["Username"] = "Benutzername darf nur Buchstaben, Zahlen, '.' und '_' enthalten."
	} else if !regex(username, "\\D") {
		(*errors)["Username"] = "Benutzername muss mindestens 1 Buchstaben enthalten."
	} else if regex(username, "^[._]") {
		(*errors)["Username"] = "Benutzername darf nicht mit '.' oder '_' beginnen."
	} else if regex(username, "[._]$") {
		(*errors)["Username"] = "Benutzername darf nicht mit '.' oder '_' enden."
	} else if regex(username, "[_.]{2}") {
		(*errors)["Username"] = "Benutzername darf '.' und '_' nicht aufeinanderfolgend haben."
	}
}

// validateEmail validates an email.
func (errors *FormErrors) validateEmail(email string) {
	if email == "" {
		(*errors)["Email"] = "Bitte Email angeben."
	} else if len(email) < 3 {
		(*errors)["Email"] = "Email muss mindestens 3 Zeichen lang sein."
	} else if len(email) > 100 {
		(*errors)["Email"] = "Email darf höchstens 100 Zeichen lang sein."
	} else if !regex(email, "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$") {
		(*errors)["Email"] = "Ungültiges Email-Format."
	}
}

// validatePassword validates a password.
func (errors *FormErrors) validatePassword(password string, errorName string) {
	if password == "" {
		(*errors)[errorName] = "Bitte Passwort angeben."
	} else if len(password) < 8 {
		(*errors)[errorName] = "Passwort muss mindestens 8 Zeichen lang sein."
	} else if !regex(password, "[!@#$%^&*]") {
		(*errors)[errorName] = "Passwort muss ein Sonderzeichen enthalten (!@#$%^&*)."
	} else if !regex(password, "[a-z]") {
		(*errors)[errorName] = "Passwort muss mindestens ein Kleinbuchstaben enthalten."
	} else if !regex(password, "[A-Z]") {
		(*errors)[errorName] = "Passwort muss mindestens ein Grossbuchstaben enthalten."
	} else if !regex(password, "\\d") {
		(*errors)[errorName] = "Passwort muss mindestens eine Zahl enthalten."
	}
}

// regex checks if a certain regular expression matches a certain string.
func regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(err)
	}

	return match
}

// TODO form validation for:
//  - phase 1 of quiz
//  - phase 2 of quiz
//  - phase 3 of quiz
