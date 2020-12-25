package web

// quiz_handler.go
// Contains all HTTP-handlers for pages evolving around playing a quiz.

import (
	"encoding/gob"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	maxIdleTimeMinutes   = 30 // max time to be spent in a specific phase of a quiz
	phase1Questions      = 3  // amount of questions in phase 1
	phase1Choices        = 3  // amount of choices per question of phase 1
	phase1Points         = 3  // amount of points per correct guess of phase 1
	phase1ChoicesMaxDiff = 10 // highest possible difference between correct year and random year of phase 1
)

// init
// Gets initialized with the package. Registers certain types to the session,
// because by default the session can only contain basic data types (int, bool,
// string, etc.).
func init() {
	gob.Register(QuizData{})
	gob.Register([]backend.Event{})
	gob.Register(backend.Event{})
}

// QuizHandler
// Object for handlers to access sessions and database.
type QuizHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// QuizData
// Contains the array of events, the user's points and the token (topic ID,
// current time and current phase) in order to validate the correct playing
// order of a quiz.
type QuizData struct {
	Events []backend.Event
	Points int

	Phase     int       // Ensures the correct playing order, so that a user can't skip any phase
	Reviewed  bool      // Ensures a user can't skip a reviewing phase
	TopicID   int       // Ensures a user can't skip from phase n of topic A to phase n of topic B
	TimeStamp time.Time // Ensures a user can't continue a quiz after n minutes of inactivity
}

// Phase1
// A GET-method that any user can call. It consists of a form with 3 multiple-
// choice questions, where the user has to guess the year of a given event.
func (handler *QuizHandler) Phase1() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		EventsMultipleChoice []phase1Question
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase1.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Benutzer eingeloggt sein, um ein Quiz zu spielen.")
			http.Redirect(res, req, "/topics/"+topicIDstr, http.StatusFound)
			return
		}

		// Convert topic ID to int
		topicID, err := strconv.Atoi(topicIDstr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Execute SQL statement to get events
		events, err := handler.store.EventsByTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Randomize array of events
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(events), func(n1, n2 int) { events[n1], events[n2] = events[n2], events[n1] })

		// Add quiz data to session
		handler.sessions.Put(req.Context(), "quiz", QuizData{
			Events:    events,
			TopicID:   topicID,
			Phase:     1,
			Reviewed:  false,
			TimeStamp: time.Now(),
		})

		// For each of the first 3 events in the array, get 2 other random
		// values for the user to guess the year from
		questions := generatePhase1Questions(events)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData:          GetSessionData(handler.sessions, req.Context()),
			EventsMultipleChoice: questions,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase1Submit
// A POST-method. It calculates the points.
func (handler *QuizHandler) Phase1Submit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicID, err := strconv.Atoi(chi.URLParam(req, "topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Retrieve quiz data from session
		quizInf := handler.sessions.Pop(req.Context(), "quiz")

		// Validate the token of the quiz data
		quiz, msg := validateQuizToken(quizInf, 1, false, topicID)
		// If msg isn't empty, an error occurred
		if msg != "" {
			handler.sessions.Put(req.Context(), "flash_error", "Ein Fehler ist aufgetreten. "+msg)
			http.Redirect(res, req, "/topics", http.StatusFound)
			return
		}

		// Loop through the 3 input forms of radio-buttons of phase 1
		for num := 0; num < phase1Questions; num++ {

			// Retrieve user's guess from form
			guess, err := strconv.Atoi(req.FormValue("q" + strconv.Itoa(num)))

			// If no error occurs, the user selected one of the three years;
			// otherwise, the user left the selection at "Ich weiss es nicht",
			// in which case we skip this question
			if err == nil {
				// Check if the user's guess is correct, by comparing it to the
				// corresponding event in the array
				if guess == quiz.Events[num].Year {
					quiz.Points += phase1Points // If guess is correct, user gets 3 points
				}
			}
		}

		// Add quiz data to session
		handler.sessions.Put(req.Context(), "quiz", quiz)

		// Redirect to review of phase 1
		http.Redirect(res, req, "/topics/"+strconv.Itoa(quiz.TopicID)+"/quiz/1/review", http.StatusFound)
	}
}

// TODO Phase1Review
// A GET-method that any user can call after Phase1. It returns a correction of the quiz.
func (handler *QuizHandler) Phase1Review() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase1.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TODO Phase2
// A GET-method that any user can call after Phase1Review. It consists of a form with 4 questions, where the
// user has to guess the year of a given event.
func (handler *QuizHandler) Phase2() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase2.html",
	))

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

// TODO Phase3
// A GET-method that any user can call after Phase2Review. It consists of a
// form with all events of the topic, where the user has to put the events in
// chronological order.
func (handler *QuizHandler) Phase3() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase3.html",
	))

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

// TODO Summary
// A GET-method that any user can call after Phase3Review. It summarizes the
// quiz played.
func (handler *QuizHandler) Summary() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_summary.html",
	))

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

// phase1Question
// Contains name of event with the correct year and 2 random years, in random
// order.
type phase1Question struct {
	Event   string // name of event
	Choices []int
	ID      string // relevant for HTML input form name
}

// generatePhase1Questions
// Generates two random numbers for each of the first three events in the array
// to use in phase 1 of the quiz (multiple-choice).
// Sample input: array of events, of which the years of the first 3 events are:
// 1945, 1960, 1981
// Sample output: [[1955 1945 1935] [1951 1961 1960] [1981 1971 1976]]
func generatePhase1Questions(events []backend.Event) []phase1Question {
	questions := make([]phase1Question, phase1Questions)

	// Set seed to generate random numbers from
	rand.Seed(time.Now().UnixNano())

	// Loop through the first events of the array
	for q := 0; q < phase1Questions; q++ {

		correctYear := events[q].Year // the event's year

		min := correctYear - phase1ChoicesMaxDiff // minimum cap of random number
		max := correctYear + phase1ChoicesMaxDiff // maximum cap of random number

		years := []int{correctYear}

		// Generate unique, random numbers between max and min, to mix with the correct year
		for c := 1; c < phase1Choices; c++ {
			rand.Seed(time.Now().Unix())          // set a seed for RNG
			newYear := rand.Intn(max-min+1) + min // generate a random number between min and max

			// Loop through array of already existing years, to check if the newly generated year is unique
			unique := true
			for _, year := range years {
				if newYear == year {
					unique = false
					break
				}
			}
			if unique {
				years = append(years, newYear) // add newly generated year to array of years
			} else {
				c-- // redo generating the previous year
			}
		}

		// Shuffle the years, so that the correct year isn't always in the
		// first spot
		rand.Shuffle(len(years), func(n1, n2 int) {
			years[n1], years[n2] = years[n2], years[n1]
		})

		// Add values to structure
		questions[q].Event = events[q].Title
		questions[q].Choices = years
		questions[q].ID = "q" + strconv.Itoa(q) // sample ID: q0
	}

	return questions
}

// validateQuizToken
// Validates the correct playing order of a quiz by comparing the phase, topic
// and time stamp of the quiz data in the session with the URL and current time
// respectively. It returns the quiz data as a structure and an empty string,
// if everything checks out, or an empty quiz data structure and an error
// string to be used in the error flash message after redirecting back.
func validateQuizToken(quizInf interface{}, phase int, reviewed bool, topicID int) (QuizData, string) {

	// Check for empty quiz interface
	if quizInf == nil {
		// Occurs when a user manually enters a URL of a later phase without
		// properly starting a quiz
		return QuizData{}, "Womöglich haben Sie versucht, unter unerlaubten Umständen ein Quiz zu starten, " +
			"ohne bei Phase 1 zu beginnen."
	}

	quiz := quizInf.(QuizData)

	// Check for invalid topic ID
	if topicID != quiz.TopicID {
		// Occurs when a user manually changes the topic ID in the URL whilst
		// in a later phase of a quiz.
		return QuizData{}, "Womöglich haben Sie versucht, während des Quizzes das Thema zu ändern."
	}

	// Check for invalid phase
	if phase != quiz.Phase || reviewed != quiz.Reviewed {
		// Occurs when a user manually changes the phase in the URL
		return QuizData{}, "Womöglich haben Sie versucht, eine Phase des Quizzes zu überspringen."
	}

	// Check for invalid time stamp. Unix() displays the time passed in seconds since
	// a specific date. By adding the time stamp of the quiz data to the max
	// idle time, we can check if it was surpassed by the current time
	if time.Now().Unix() > quiz.TimeStamp.Unix()+maxIdleTimeMinutes*60 {
		// Occurs when a user takes unreasonably long to complete a phase
		return QuizData{}, "Nach " + strconv.Itoa(maxIdleTimeMinutes) + " Minuten Inaktivität in einer Phase, " +
			"endet das Quiz, da angenommen wird, der Benutzer habe das Quiz verlassen."
	}

	return quiz, ""
}
