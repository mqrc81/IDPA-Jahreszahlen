package web

// play_handler.go
// Contains all HTTP-handlers for pages evolving around playing a game.

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// PlayHandler
// Object for handlers to access sessions and database.
type PlayHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// Phase1
// A GET-method that any user can call. It consists of a form with 3 multiple-
// choice questions, where the user has to guess the year of a given event.
func (h *PlayHandler) Phase1() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData

		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/play_phase1.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, err := strconv.Atoi(topicIDstr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get events
		ee, err := h.store.EventsByTopic(topicID, true)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO session management to ensure correct play order etc.

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Events:      ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase2
// A GET-method that any user can call after having completed Phase1. It consists of a form with 4 questions, where the
// user has to guess the year of a given event.
func (h *PlayHandler) Phase2() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData

		Events []backend.Event
	}

	// Parse HTML-template
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/play_phase2.html"))

	return func(res http.ResponseWriter, req *http.Request) {

		// TODO
		// Retrieve values from URL parameters
		// topicIDstr := chi.URLParam(req, "topicID")
		// topicID, err := strconv.Atoi(topicIDstr)
		// if err != nil {
		//	http.Error(res, err.Error(), http.StatusInternalServerError)
		//	return
		// }

		// Check if game order is legal
		// check := session.PopString("game")
		// if check == nil || strings.Split(check, ":")[0] != "1" || strings.Split(check, ":")[1] != topicIDstr {
		//	http.Error(res, "Illegal game session", http.StatusBadRequest)
		//	return
		// }

		// TODO Retrieve events array from sessions
		var ee []backend.Event // TEMP

		points := 0
		for x := 1; x <= 3; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q" + strconv.Itoa(x)))
			if guess == ee[x].Year {
				points += 3
			}
		}
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Events:      ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Phase3
// A GET-method that any user can call after having completed Phase2. It consists of a form with up to 15 questions,
// where the user has to match the year to any of the given events.
func (h *PlayHandler) Phase3() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData

		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/play_phase3.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from session
		var ee []backend.Event // ee := session.Pop(ctx, "events")
		var points int         // points := session.PopInt(ctx, "points")

		for x := 1; x <= 5; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q" + strconv.Itoa(x))) // The user's answer
			correctYear := ee[x+3].Year                                    // The correct answer

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
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		// Execute SQL statement to get events
		ee, err = h.store.EventsByTopic(topicID, true)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add values to session
		// session.Put(ctx, "points", points)
		// session.Put(ctx, "events", ee)

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Events:      ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Store
// A POST-method. It stores score of game played and redirects to Review.
func (h *PlayHandler) Store() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// TODO Retrieve values from session
		var userID int // Temporary
		// var ee []backend.Event // Temporary
		var points int // Temporary

		for x := 1; x <= 3; x++ {
			// TODO algorithm for points in phase 3
		}

		if err := h.store.CreateScore(&backend.Score{
			ScoreID: 0,
			TopicID: topicID,
			UserID:  userID,
			Points:  points,
			Date:    "",
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		// Redirects to review of the user's score
		http.Redirect(res, req, "/topics/"+topicIDstr+"/play/review", http.StatusFound)
	}
}

// Review
// A GET-method that any user can call after having completed Phase3. It
// summarizes the game played.
func (h *PlayHandler) Review() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
		// TODO
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/play_review.html"))

	return func(res http.ResponseWriter, req *http.Request) {

		// TODO

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {

		}
	}
}
