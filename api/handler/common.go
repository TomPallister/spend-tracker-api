package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TomPallister/godutch-api/api/infrastructure"
	"github.com/gorilla/mux"
)

// ErrorNotFound ...to be used when object does not exist in repository
var ErrorNotFound = errors.New("Tracker not found")

// GetIDFromVARs ...
func GetIDFromVARs(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetCryptoEmailFromVars ...
func GetCryptoEmailFromVars(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	cryptoEmail, ok := vars["cryptoEmail"]
	if ok == false {
		return "", ErrorNotFound
	}
	return cryptoEmail, nil
}

// CreateErrorResponseAndLog ...
func CreateErrorResponseAndLog(headerValue int, w http.ResponseWriter, logger infrastructure.Logger, err error) {
	w.WriteHeader(headerValue)
	createErrorJSONResponse(err, w)
	logger.Error("There was an error", err)
	return
}

type errorResponse struct {
	Error string
}

func createErrorJSONResponse(err error, w http.ResponseWriter) {
	response := errorResponse{Error: err.Error()}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
