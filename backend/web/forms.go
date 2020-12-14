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
// String, etc.).
func init() {
	gob.Register(TopicForm{})
	gob.Register(EventForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
	gob.Register(PasswordForm{})
	gob.Register(FormErrors{})
}

// FormErrors
// A map that holds the error messages.
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
	if f.Year <= 0 {
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

	// Validate username
	if f.UsernameTaken {
		f.Errors["Username"] = "Benutzername ist bereits vergeben."
	} else if len(f.Username) < 3 {
		f.Errors["Username"] = "Benutzername muss mindestens 3 Zeichen lang sein."
	} else if len(f.Username) > 20 {
		f.Errors["Username"] = "Benutzername darf höchstens 20 Zeichen lang sein."
	} else if !Regex(f.Username, "^[a-zA-Z0-9._]*$") {
		f.Errors["Username"] = "Benutzername darf nur Buchstaben, Zahlen, '.' und '_' enthalten."
	} else if !Regex(f.Username, "\\D") {
		f.Errors["Username"] = "Benutzername muss mindestens 1 Buchstaben enthalten."
	} else if Regex(f.Username, "^[._]") {
		f.Errors["Username"] = "Benutzername darf nicht mit '.' oder '_' beginnen."
	} else if Regex(f.Username, "[._]$") {
		f.Errors["Username"] = "Benutzername darf nicht mit '.' oder '_' enden."
	} else if Regex(f.Username, "[_.]{2}") {
		f.Errors["Username"] = "Benutzername darf '.' und '_' nicht aufeinanderfolgend haben."
	}

	// Validate password
	f.Errors = ValidatePassword(f.Password, f.Errors, "Password")

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

// PasswordForm
// Holds values of the form input when editing a password.
type PasswordForm struct {
	NewPassword          string
	IncorrectOldPassword bool

	Errors FormErrors
}

// Validate
// Validates the form input when editing a password.
func (f *PasswordForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate old password
	if f.IncorrectOldPassword {
		f.Errors["OldPassword"] = "Altes Passwort ist inkorrekt."
	}

	// Validate new password
	f.Errors = ValidatePassword(f.NewPassword, f.Errors, "NewPassword")

	return len(f.Errors) == 0
}

// Regex
// Checks if a certain regular expression matches a certain string
func Regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

// ValidatePassword
// Validates a password.
func ValidatePassword(password string, formErrors FormErrors, errorName string) FormErrors {
	if len(password) < 8 {
		formErrors[errorName] = "Passwort muss mindestens 8 Zeichen lang sein."
	} else if !Regex(password, "[!@#$%^&*]") {
		formErrors[errorName] = "Passwort muss ein Sonderzeichen enthalten (!@#$%^&*)."
	} else if !Regex(password, "[a-z]") {
		formErrors[errorName] = "Passwort muss mindestens ein Kleinbuchstaben enthalten."
	} else if !Regex(password, "[A-Z]") {
		formErrors[errorName] = "Passwort muss mindestens ein Grossbuchstaben enthalten."
	} else if !Regex(password, "\\d") {
		formErrors[errorName] = "Passwort muss mindestens eine Zahl enthalten."
	}

	return formErrors
}
