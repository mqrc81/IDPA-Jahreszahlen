package web

/*
 * Contains all HTTP-handler functions for pages evolving around playing a
 * quiz.
 */

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
	timeExpiry = 20 // max time to be spent in a specific phase of a quiz

	p1Questions      = 3  // amount of questions in phase 1
	p1Choices        = 3  // amount of choices per question of phase 1
	p1Points         = 3  // amount of points per correct guess of phase 1
	p1ChoicesMaxDiff = 10 // highest possible difference between the correct year and a random year of phase 1

	p2Questions     = 4 // amount of questions in phase 2
	p2Points        = 8 // amount of points per correct guess of phase 2
	p2PartialPoints = 3 // amount of partial points possible in phase 2, when guess was incorrect, but close
)

// init gets initialized with the package. It registers certain types to the
// session, because by default the session can only contain basic data types
// (int, bool, string, etc.).
func init() {
	gob.Register(QuizData{})
	gob.Register(backend.Topic{})
	gob.Register([]backend.Event{})
	gob.Register(backend.Event{})
	gob.Register([]phase1Question{})
	gob.Register(phase1Question{})
	gob.Register([]int{})
	gob.Register([]phase2Question{})
	gob.Register(phase2Question{})
}

// QuizHandler is the object for handlers to access sessions and database.
type QuizHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// QuizData contains the topic with the array of events and the points to keep
// track of, as well as the equivalent of a token (consisting of the topic ID,
// time expiry and current phase) in order to validate the correct playing
// order of a quiz.
type QuizData struct {
	Topic  backend.Topic // Contains topic ID for validation and events for playing the quiz
	Points int

	Phase     int       // Ensures the correct playing order, so that a user can't skip any phase
	Reviewed  bool      // Ensures a user can't skip a reviewing phase
	TimeStamp time.Time // Ensures a user can't continue a quiz after n minutes of inactivity
}

// Phase1 is a GET-method that any user can call. It consists of a form with 3
// multiple-choice questions, where the user has to guess the year of a given
// event.
func (handler *QuizHandler) Phase1() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Questions []phase1Question
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase1.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, err := strconv.Atoi(topicIDstr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Benutzer eingeloggt sein, um ein Quiz zu spielen.")
			http.Redirect(res, req, "/topics/"+topicIDstr, http.StatusFound)
			return
		}

		// Execute SQL statement to get topic
		topic, err := handler.store.GetTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Shuffle array of events
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(topic.Events), func(n1, n2 int) {
			topic.Events[n1], topic.Events[n2] = topic.Events[n2], topic.Events[n1]
		})

		// Add quiz data to session for future phases
		handler.sessions.Put(req.Context(), "quiz", QuizData{
			Topic:     topic,
			Phase:     1,
			Reviewed:  false,
			TimeStamp: time.Now(),
		})

		// For each of the first 3 events in the array, generate 2 other random
		// years for the user to guess from
		questions := generatePhase1Questions(topic.Events)

		// Add questions to session for review of phase 1
		handler.sessions.Put(req.Context(), "phase1questions", questions)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Questions:   questions,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase1Submit is a POST-method. It calculates the points and redirects to
// Phase1Review.
func (handler *QuizHandler) Phase1Submit() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")

		// Retrieve quiz data from session
		quiz := handler.sessions.Get(req.Context(), "quiz").(QuizData)

		// Loop through the 3 input forms of radio-buttons of phase 1
		for num := 0; num < p1Questions; num++ {

			// Retrieve user's guess from form
			guess, _ := strconv.Atoi(req.FormValue("q" + strconv.Itoa(num)))

			// Check if the user's guess is correct, by comparing it to the
			// corresponding event in the array of events of the topic
			if guess == quiz.Topic.Events[num].Year {
				quiz.Points += p1Points // If guess is correct, user gets 3 points
			}
		}

		// Redirect to review of phase 1 whilst keeping the previous request
		// context
		http.Redirect(res, req.WithContext(req.Context()), "/topics/"+topicIDstr+"/quiz/1/review", http.StatusFound)
	}
}

// Phase1Review is a GET-method that any user can call after Phase1. It
// displays a correction of the questions.
func (handler *QuizHandler) Phase1Review() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Questions []phase1Question
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase1_review.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve quiz data from session
		quizInf := handler.sessions.Get(req.Context(), "quiz")

		// Validate the token of the quiz data
		quiz, msg := validateQuizToken(quizInf, 1, false, topicID)
		// If msg isn't empty, an error occurred
		if msg != "" {
			handler.sessions.Put(req.Context(), "flash_error",
				"Ein Fehler ist aufgetreten in Phase 1 eines Quizzes. "+msg)
			http.Redirect(res, req, "/topics", http.StatusFound)
			return
		}

		// Update quiz data
		quiz.Reviewed = true
		quiz.TimeStamp = time.Now()

		// Retrieve questions from session
		questions := handler.sessions.Get(req.Context(), "phase1questions").([]phase1Question)

		// Add quiz data to session for later phases
		handler.sessions.Put(req.Context(), "quiz", quiz)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Questions:   questions,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase2 is a GET-method that any user can call after Phase1Review. It
// consists of a form with 4 questions, where the user has to guess the year of
// a given event.
func (handler *QuizHandler) Phase2() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Questions []phase2Question
	}

	// Parse HTML-templates
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

// TODO
// Phase3 is a GET-method that any user can call after Phase2Review. It
// consists of a form with all events of the topic, where the user has to put
// the events in chronological order.
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

// TODO
// Summary is a GET-method that any user can call after Phase3Review. It
// summarizes the quiz completed.
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
	if topicID != quiz.Topic.TopicID {
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
	if time.Now().Unix() > quiz.TimeStamp.Unix()+timeExpiry*60 {
		// Occurs when a user takes unreasonably long to complete a phase
		return QuizData{}, "Nach " + strconv.Itoa(timeExpiry) + " Minuten Inaktivität in einer Phase, " +
			"endet das Quiz, da angenommen wird, der Benutzer habe das Quiz verlassen."
	}

	return quiz, ""
}

// phase1Question represents 1 of the 3 multiple-choice questions of phase 1.
// It contains name of event, year of event and 2 random years randomly mixed
// in with the correct year.
type phase1Question struct {
	EventName string // name of event
	EventYear int    // year of event
	Choices   []int  // choices in random order (including correct year)
	ID        string // only relevant for HTML input form name
}

// generatePhase1Questions generates 3 phase1Question structures by generating
// 2 random years for each of the first 3 events in the array.
//
// Sample input: []backend.Event{{..., Year: 1945}, {..., Year: 1960}, {..., Year: 1981}, ...}
// Sample output: [[1955 1945 1938] [1951 1961 1960] [1981 1971 1976]]
func generatePhase1Questions(events []backend.Event) []phase1Question {
	// Create non-nil array of questions
	questions := make([]phase1Question, p1Questions)

	// Set seed to generate random numbers from
	rand.Seed(time.Now().UnixNano())

	// Loop through the first events of the array
	for q := 0; q < p1Questions; q++ {

		correctYear := events[q].Year // the event's year

		min := correctYear - p1ChoicesMaxDiff // minimum cap of random number
		max := correctYear + p1ChoicesMaxDiff // maximum cap of random number

		years := []int{correctYear}

		// Generate unique, random numbers between max and min, to mix with the correct year
		for c := 1; c < p1Choices; c++ {
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
		questions[q].EventName = events[q].Name
		questions[q].EventYear = events[q].Year
		questions[q].Choices = years
		questions[q].ID = "q" + strconv.Itoa(q) // sample ID: q0
	}

	return questions
}

// phase2Question represents 1 of the 4 questions of phase 2. It contains name
// of event and year of event.
type phase2Question struct {
	EventName string // name of event
	EventYear int    // year of event
	ID        string // only relevant for HTML input form name
}

// TODO
// generatePhase2Questions generates 4 phase2Question structures from events 4-8
// of the topic .
func generatePhase2Questions(events []backend.Event) []phase2Question {
	// Create non-nil array of questions
	questions := make([]phase2Question, p2Questions)

	return questions
}
