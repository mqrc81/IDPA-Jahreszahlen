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
	idleTimeMinutes = 15
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

	Phase    int       // Ensures the correct playing order, so that a user can't skip any phase
	Reviewed bool      // Ensures a user can't skip a reviewing phase
	TopicID  int       // Ensures a user can't skip from phase n of topic A to phase n of topic B
	Time     time.Time // Ensures a user can't continue a quiz after n minutes of inactivity
}

// Phase1
// A GET-method that any user can call. It consists of a form with 3 multiple-
// choice questions, where the user has to guess the year of a given event.
func (handler *QuizHandler) Phase1() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		EventsMultipleChoice []eventMultipleChoice
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
			Events:   events,
			TopicID:  topicID,
			Phase:    1,
			Reviewed: false,
			Time:     time.Now(),
		})

		// For each of the 3 first events in the array, get two other random
		// values for the user to guess the year from
		triplets := generateEventMultipleChoice(events)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData:          GetSessionData(handler.sessions, req.Context()),
			EventsMultipleChoice: triplets,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase1Submit
// A POST-method. It validates the form from Phase1 and calculates the points.
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
		quiz, msg := validateToken(quizInf, 1, false, topicID)
		// If msg isn't empty, an error occurred
		if msg != "" {
			handler.sessions.Put(req.Context(), "flash_error", "Ein Fehler ist aufgetreten. "+msg)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Loop through the 3 input forms of radio-buttons of phase 1
		for num := 1; num <= 3; num++ {

			// Retrieve user's guess from form
			guess, err := strconv.Atoi(req.FormValue("q" + strconv.Itoa(num)))

			// If no error occurs, the user selected one of the three years.
			// Otherwise, the user left the selection at "Ich weiss es nicht",
			// which means we skip this question.
			if err == nil {

				// Check if the user's guess is correct, by comparing it to the
				// corresponding event in the array
				if guess == quiz.Events[num-1].Year {
					quiz.Points += 3 // If answer is correct, user gets 3 points
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

		Events []backend.Event
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase2.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// TODO
		// Retrieve values from URL parameters
		// topicIDstr := chi.URLParam(req, "topicID")
		// topicID, err := strconv.Atoi(topicIDstr)
		// if err != nil {
		//	http.Error(res, err.Error(), http.StatusInternalServerError)
		//	return
		// }

		var ee []backend.Event // TEMP

		points := 0
		for x := 1; x <= 3; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q" + strconv.Itoa(x)))
			if guess == ee[x].Year {
				points += 3
			}
		}
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Events:      ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TODO Phase3
// A GET-method that any user can call after having completed Phase2. It consists of a form with up to 15 questions,
// where the user has to match the year to any of the given events.
func (handler *QuizHandler) Phase3() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_phase3.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from session
		var events []backend.Event // events := session.Pop(ctx, "events")
		var points int             // points := session.PopInt(ctx, "points")

		for x := 1; x <= 5; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q" + strconv.Itoa(x))) // The user's answer
			correctYear := events[x+3].Year                                // The correct answer

			// If answer is correct, user gets 7 points
			if guess == correctYear {
				points += 7
			} else {
				// If answer is close, he gets partial points
				diff := 0
				if guess > correctYear {
					diff = guess - correctYear
				} else if correctYear > guess {
					diff = correctYear - guess
				}
				points += 4 - diff
			}
		}

		// Retrieve topic ID from URL
		topicID, err := strconv.Atoi(req.URL.Query().Get("topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
		}

		// Execute SQL statement to get events
		events, err = handler.store.EventsByTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Randomize array of events
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(events), func(n1, n2 int) { events[n1], events[n2] = events[n2], events[n1] })

		// Add values to session
		// session.Put(ctx, "points", points)
		// session.Put(ctx, "events", events)

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Events:      events,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TODO Store
// A POST-method. It stores score of quiz played and redirects to Summary.
func (handler *QuizHandler) Store() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// TODO Retrieve values from session
		var userID int // Temporary
		// var ee []backend.Event // Temporary
		var points int // Temporary

		for x := 1; x <= 3; x++ {
			// TODO algorithm for points in phase 3
		}

		if err := handler.store.CreateScore(&backend.Score{
			TopicID: topicID,
			UserID:  userID,
			Points:  points,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		// Redirects to review of the user's score
		http.Redirect(res, req, "/topics/"+topicIDstr+"/quiz/overview", http.StatusFound)
	}
}

// TODO Summary
// A GET-method that any user can call after Phase3Review. It summarizes the
// quiz played.
func (handler *QuizHandler) Summary() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
		// TODO
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/quiz_summary.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// TODO

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// eventMultipleChoice
// Contains name of event with the correct year and 2 random years, in random
// order.
type eventMultipleChoice struct {
	Event   string
	Triplet [3]int
}

// generateEventMultipleChoice
// Generates two random numbers for each of the first three events in the array
// to use in phase 1 of the quiz (multiple-choice).
// Sample input: array of events, of which the years of the first 3 events are:
// 1945, 1960, 1981
// Sample output: [[1955 1945 1935] [1951 1961 1960] [1981 1971 1976]]
func generateEventMultipleChoice(events []backend.Event) []eventMultipleChoice {
	eventsTriplets := make([]eventMultipleChoice, 3)

	// Set seed to generate random numbers from
	rand.Seed(time.Now().UnixNano())

	// Loop through the 3 first events of the array
	for x := 0; x < 3; x++ {

		correctYear := events[x].Year // the event's year
		randomYear1 := correctYear    // random number #1
		randomYear2 := correctYear    // random number #2

		min := correctYear - 10 // minimum cap of random number
		max := correctYear + 10 // maximum cap of random number

		// Generate a unique random number between max and min
		for randomYear1 == correctYear {
			randomYear1 = rand.Intn(max-min+1) + min
		}
		for randomYear2 == correctYear || randomYear2 == randomYear1 {
			randomYear2 = rand.Intn(max-min+1) + min
		}

		// Put triplet of years into a temporary array
		temp := [3]int{correctYear, randomYear1, randomYear2}

		// Shuffle the years, so that the correct year isn't always in the
		// first spot
		rand.Shuffle(len(temp), func(i, j int) {
			temp[i], temp[j] = temp[j], temp[i]
		})

		// Add values to struct
		eventsTriplets[x].Event = events[x].Title
		eventsTriplets[x].Triplet = temp
	}

	return eventsTriplets
}

// validateToken
// Validates the correct playing order of a quiz by comparing the phase, topic
// and time stamp of the quiz data in the session with the URL and current time
// respectively. It returns the quiz data as a struct and an empty string, if
// everything checks out, or an empty quiz data struct and an error string to
// be used in the error flash message after redirecting back.
func validateToken(quiz interface{}, phase int, reviewed bool, topicID int) (QuizData, string) {
	// TODO check for all token error cases
	return QuizData{}, ""
}
