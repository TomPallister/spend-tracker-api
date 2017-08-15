package spendsummarieshandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TomPallister/godutch-api/api/environment"
	"github.com/TomPallister/godutch-api/api/handler"
	"github.com/TomPallister/godutch-api/api/view"
)

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

		spendSummaries, err := env.SpendSummaryService.FindSpendSummariesForTrackerID(subject, id)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return 
		}

		spendSummariesView := view.SpendSummaries{
			SpendSummaries: spendSummaries,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(spendSummariesView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}
