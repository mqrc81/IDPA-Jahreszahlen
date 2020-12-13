package web

import (
	"encoding/gob"
	"log"
	"regexp"
	"time"
)

/*
 * init gets called the first time this file is initialized
 */
func init() {
	gob.Register(TopicForm{})
	gob.Register(CreateEventForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
	gob.Register(PasswordForm{})
	gob.Register(FormErrors{})
}

// FormErrors is a map that holds the errors
type FormErrors map[string]string

// TopicForm holds values of form when creating a topic
type TopicForm struct {
	Title       string
	StartYear   int
	EndYear     int
	Description string

	Errors FormErrors
}

/*
 * Validate validates topic form
 */
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

// CreateEventForm holds values of form when creating a event
type CreateEventForm struct {
	Title string
	Year  int

	Errors FormErrors
}

/*
 * Validate validates event form
 */
func (f *CreateEventForm) Validate() bool {
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

// RegisterForm holds values of form when registering a user
type RegisterForm struct {
	Username      string
	Password      string
	UsernameTaken bool

	Errors FormErrors
}

/*
 * Regex checks if the regular expression matches the string
 */
func Regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

/*
 * ValidatePassword validates a user's password
 */
func ValidatePassword(password string, errors FormErrors, errorName string) FormErrors {
	if len(password) < 8 {
		errors[errorName] = "Passwort muss mindestens 8 Zeichen lang sein."
	} else if !Regex(password, "[!@#$%^&*]") {
		errors[errorName] = "Passwort muss ein Sonderzeichen enthalten (!@#$%^&*)."
	} else if !Regex(password, "[a-z]") {
		errors[errorName] = "Passwort muss mindestens ein Kleinbuchstaben enthalten."
	} else if !Regex(password, "[A-Z]") {
		errors[errorName] = "Passwort muss mindestens ein Grossbuchstaben enthalten."
	} else if !Regex(password, "\\d") {
		errors[errorName] = "Passwort muss mindestens eine Zahl enthalten."
	}

	return errors
}

/*
 * Validate validates form
 */
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

// LoginForm holds values of form when registering a user
type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool

	Errors FormErrors
}

/*
 * Validate validates form
 */
func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate username
	if f.Username == "" {
		f.Errors["Username"] = "Bitte Benutzernamen eingeben."
	} else if f.IncorrectCredentials {
		f.Errors["Username"] = "Benutzername oder Passwort ist falsch."
	}

	// Validate password
	if f.Password == "" {
		f.Errors["Password"] = "Bitte Passwort eingeben."
	} else if f.IncorrectCredentials {
		f.Errors["Username"] = " "
	}

	return len(f.Errors) == 0
}

// PasswordForm holds values of form when editing a user's password
type PasswordForm struct {
	Password1            string
	Password2            string
	IncorrectOldPassword bool

	Errors FormErrors
}

/*
 * Validate validates form
 */
func (f *PasswordForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate old password
	if f.IncorrectOldPassword {
		f.Errors["OldPassword"] = "Altes Passwort ist inkorrekt."
	}

	// Validate new password
	f.Errors = ValidatePassword(f.Password1, f.Errors, "Password1")

	// Validate new password confirmed
	f.Errors = ValidatePassword(f.Password2, f.Errors, "Password2")

	return len(f.Errors) == 0
}
