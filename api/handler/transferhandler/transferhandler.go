package transferhandler

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

		trackerID := r.URL.Query()["trackerId"] 
		id, err := strconv.ParseInt(trackerID[0], 10, 64)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
 
		transfers, err := env.TransferService.FindTransfersForTrackerID(subject, id)
		if err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}

		transfersView := view.Transfers{
			Transfers: transfers,
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(transfersView); err != nil {
			handler.CreateErrorResponseAndLog(http.StatusBadRequest, w, env.Logger, err)
			return
		}
	})
}
