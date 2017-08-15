package handler

import (
	"encoding/json"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	respondJson("All good. You don't need to be authenticated to call this", w)
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	respondJson("All good. You only get this message if you're authenticated", w)
}

type Response struct {
	Text string `json:"text"`
}

func respondJson(text string, w http.ResponseWriter) {
	response := Response{text}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
