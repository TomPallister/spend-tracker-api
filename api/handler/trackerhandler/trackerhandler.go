package trackerhandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/handler"
	"github.com/TomPallister/godutch-api/api/view"
)

// CreateTrackerHandler ...
func CreateTrackerHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		var tracker view.Tracker
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &tracker); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		newTracker, err := env.TrackerService.CreateTracker(subject, tracker.Tracker)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		trackerView := view.Tracker{
			Tracker: newTracker,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(trackerView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// UpdateTrackerHandler ...
func UpdateTrackerHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		id, err := handler.GetIDFromVARs(r)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		var tracker view.Tracker
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &tracker); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if id == 0 {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		tracker.Tracker.ID = id

		updatedTracker, err := env.TrackerService.UpdateTracker(subject, tracker.Tracker)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		trackerView := view.Tracker{
			Tracker: updatedTracker,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(trackerView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// DeleteTrackerHandler ...
func DeleteTrackerHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		id, err := handler.GetIDFromVARs(r)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		result, err := env.TrackerService.DeleteTracker(subject, id)
		if err != nil || result == false {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNoContent)

	})
}

// FindTrackersBySubHandler ...
func FindTrackersBySubHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		trackers, err := env.TrackerService.FindByUser(subject)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		trackersView := view.Trackers{
			Trackers: trackers,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(trackersView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// FindTrackersByIDHandler ...
func FindTrackersByIDHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		id, err := handler.GetIDFromVARs(r)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		tracker, err := env.TrackerService.FindByID(subject, id)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		trackerView := view.Tracker{
			Tracker: tracker,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(trackerView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}
