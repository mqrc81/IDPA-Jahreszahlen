package web

/*
 * play_handler.go contains HTTP-handler functions for games
 */

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type PlayHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

/*
 * Phase1 is a GET-method with a form with 3 multiple-choice questions
 */
func (h *PlayHandler) Phase1() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Events []backend.Event
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(`TODO`)) // TODO

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, err := strconv.Atoi(topicIDstr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement
		ee, err := h.store.EventsByTopic(topicID, true)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Ensure correct play order
		// session.Put(ctx, "game", "1:"+topicIDstr)

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			Events: ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Phase2 is a GET-method with a form with 3 questions
 */
func (h *PlayHandler) Phase2() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Events []backend.Event
	}

	//Parse HTML-template
	tmpl := template.Must(template.New("").Parse(`TODO`)) // TODO
	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from URL parameters
		//topicIDstr := chi.URLParam(req, "topicID")
		//topicID, err := strconv.Atoi(topicIDstr)
		//if err != nil {
		//	http.Error(res, err.Error(), http.StatusInternalServerError)
		//	return
		//}

		// Check if game order is legal
		// check := session.PopString("game")
		//if check == nil || strings.Split(check, ":")[0] != "1" || strings.Split(check, ":")[1] != topicIDstr {
		//	http.Error(res, "Illegal game session", http.StatusBadRequest)
		//	return
		//}

		// TODO Retrieve events array from sessions
		var ee []backend.Event // Temporary

		points := 0
		for x := 1; x <= 3; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q"+strconv.Itoa(x)))
			if guess == ee[x].Year {
				points += 3
			}
		}
		if err := tmpl.Execute(res, data{
			Events: ee,
		}); err != nil {
		    http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Phase3 is a GET-method with a form with timeline questions
 */
func (h *PlayHandler) Phase3() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Events []backend.Event
	}

	//Parse HTML-templat
	tmpl := template.Must(template.New("").Parse(`TODO`)) // TODO
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from session
		var ee []backend.Event // ee := session.Pop(ctx, "events")
		var points int // points := session.PopInt(ctx, "points")

		for x := 1; x <= 5; x++ {
			guess, _ := strconv.Atoi(req.FormValue("q"+strconv.Itoa(x))) // The user's answer
			yr := ee[x+3].Year // The actual answer

			// If answer is correct, user gets 7 points
			if guess == yr {
				points += 7
			} else {
				// If answer is close, he gets partial points
				diff := 0
				if guess > yr {
					diff = guess - yr
				} else if yr > guess {
					diff = yr - guess
				}
				points += 4 - diff
			}
		}

		// Execute SQL statement
		topicID, err := strconv.Atoi(req.URL.Query().Get("topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
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
			Events: ee,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Submit is a POST-method that stores topic, user, points and date in database
 */
func (h *PlayHandler) Submit() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// TODO Retrieve values from session
		var userID int // Temporary
		//var ee []backend.Event // Temporary
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
		http.Redirect(res, req, "/topics/"+topicIDstr+"/play/review", http.StatusFound)
	}
}

/*
 * Review tells player his score
 */
func (h *PlayHandler) Review() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		
	}
	return func(res http.ResponseWriter, req *http.Request) {

		// TODO

	}
}
