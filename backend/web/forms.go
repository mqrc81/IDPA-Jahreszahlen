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
	switch {
	case f.Title == "":
		f.Errors["Title"] = "Titel darf nicht leer sein."
	case len(f.Title) > 50:
		f.Errors["Title"] = "Titel darf 50 Zeichen nicht überschreiten."
	}

	// Validate start- and end-year
	current := time.Now().Year()
	switch {
	case f.StartYear <= 0:
		f.Errors["Year"] = "Start-Jahr muss positiv sein."
	case f.EndYear <= 0:
		f.Errors["Year"] = "End-Jahr muss positiv sein."
	case f.StartYear > current:
		f.Errors["Year"] = "Start-Jahr darf nicht in der Zukunft sein."
	case f.EndYear > current:
		f.Errors["Year"] = "End-Jahr darf nicht in der Zukunft sein."
	case f.EndYear < f.StartYear:
		f.Errors["Year"] = "Da wurden wohl Start- und End-Jahr vertauscht."
	}

	// Validate description
	switch {
	case len(f.Description) > 500:
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
	switch {
	case f.Title == "":
		f.Errors["Title"] = "Titel darf nicht leer sein."
	case len(f.Title) > 110:
		f.Errors["Title"] = "Titel darf 110 Zeichen nicht überschreiten."
	}

	// Validate year
	switch {
	case f.Year <= 0:
		f.Errors["Year"] = "Jahr muss positiv sein."
	case f.Year > time.Now().Year():
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

func Regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

/*
 * Validate validates form
 */
func (f *RegisterForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate username
	switch {
	case f.UsernameTaken:
		f.Errors["Username"] = "Benutzername ist bereits vergeben."
	case len(f.Username) < 3:
		f.Errors["Username"] = "Benutzername muss mindestens 3 Zeichen lang sein."
	case len(f.Username) > 20:
		f.Errors["Username"] = "Benutzername darf höchstens 20 Zeichen lang sein."
	case !Regex(f.Username, "^[a-zA-Z0-9._]*$"):
		f.Errors["Username"] = "Benutzername darf nur Buchstaben, Zahlen, '.' und '_' enthalten."
	case !Regex(f.Username, "\\D"):
		f.Errors["Username"] = "Benutzername muss mindestens 1 Buchstaben enthalten."
	case Regex(f.Username, "^[._]"):
		f.Errors["Username"] = "Benutzername darf nicht mit '.' oder '_' beginnen."
	case Regex(f.Username, "[._]$"):
		f.Errors["Username"] = "Benutzername darf nicht mit '.' oder '_' enden."
	case Regex(f.Username, "[_.]{2}"):
		f.Errors["Username"] = "Benutzername darf '.' und '_' nicht aufeinanderfolgend haben."
	}

	// Validate password
	switch {
	case len(f.Password) < 8:
		f.Errors["Password"] = "Passwort muss mindestens 8 Zeichen lang sein."
	case !Regex(f.Password, "[!@#$%^&*]"):
		f.Errors["Password"] = "Passwort muss ein Sonderzeichen enthalten (!@#$%^&*)."
	case !Regex(f.Password, "[a-z]"):
		f.Errors["Password"] = "Passwort muss mindestens ein Kleinbuchstaben enthalten."
	case !Regex(f.Password, "[A-Z]"):
		f.Errors["Password"] = "Passwort muss mindestens ein Grossbuchstaben enthalten."
	case !Regex(f.Password, "\\d"):
		f.Errors["Password"] = "Passwort muss mindestens eine Zahl enthalten."
	}

	return len(f.Errors) == 0
}

// LoginForm holds values of form when registering a user
type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool

	Errors FormErrors
}

func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}

	// Validate username
	switch {
	case f.Username == "":
		f.Errors["Username"] = "Bitte Benutzernamen eingeben."
	case f.IncorrectCredentials:
		f.Errors["Username"] = "Benutzername oder Passwort ist falsch."
	}
	switch {
	case f.Password == "":
		f.Errors["Password"] = "Bitte Passwort eingeben."
	case f.IncorrectCredentials:
		f.Errors["Username"] = " "
	}

	return len(f.Errors) == 0
}
