package spendhandler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/handler"
	"github.com/TomPallister/godutch-api/api/view"
)

// CreateSpendHandler ...
func CreateSpendHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		var spend view.Spend
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &spend); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		newSpend, err := env.SpendService.CreateSpend(subject, spend.Spend)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		viewSpend := view.Spend{
			Spend: newSpend,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(viewSpend); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// UpdateSpendHandler ...
func UpdateSpendHandler(env *environment.Env) http.Handler {
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

		var spend view.Spend
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := r.Body.Close(); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if err := json.Unmarshal(body, &spend); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		if id == 0 {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		spend.Spend.ID = id

		updatedSpend, err := env.SpendService.UpdateSpend(subject, spend.Spend)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		viewSpend := view.Spend{
			Spend: updatedSpend,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(viewSpend); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}

// DeleteSpendHandler ...
func DeleteSpendHandler(env *environment.Env) http.Handler {
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

		result, err := env.SpendService.DeleteSpend(subject, id)
		if err != nil && result == true {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNoContent)
	})
}

// FindFindByTrackerIDHandler ...
func FindFindByTrackerIDHandler(env *environment.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		subject, err := env.SubjectFinder.FindSubject(r, env.Logger)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusUnauthorized, w, env.Logger, err)
			return
		}

		id, err := strconv.ParseInt(r.URL.Query()["trackerId"][0], 10, 64)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		spends, err := env.SpendService.FindByTrackerID(subject, id)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		viewSpends := view.Spends{
			Spends: spends,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(viewSpends); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}
