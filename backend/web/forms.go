package web

import (
	"encoding/gob"
	"fmt"
	"time"
)

/*
 * init gets called the first time this file is initialized
 */
func init() {
	gob.Register(CreateTopicForm{})
	gob.Register(CreateEventForm{})
	gob.Register(FormErrors{})
}

// FormErrors is a map that holds the errors
type FormErrors map[string]string

// CreateTopicForm holds values of form when creating a topic
type CreateTopicForm struct {
	Title       string
	StartYear   int
	EndYear     int
	Description string

	Errors FormErrors
}

/*
 * Validate validates form
 */
func (f *CreateTopicForm) Validate() bool {
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

	fmt.Print(f.Errors)
	return len(f.Errors) == 0
}

// CreateEventForm holds values of form when creating a event
type CreateEventForm struct {
	Title string
	Year  int

	Errors FormErrors
}

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

	fmt.Print(f.Errors)
	return len(f.Errors) == 0
}
