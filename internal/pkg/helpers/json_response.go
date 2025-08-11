package helpers

import (
	"encoding/json"
	"net/http"
)

type ResponseWrapper struct {
	Info *Info `json:"info,omitzero"`
	Data any   `json:"data,omitzero"`
}

type Info struct {
	TotalRecords int `json:"totalRecords"`
	CurrentPage  int `json:"currentPage,omitzero"`
	PageSize     int `json:"pageSize,omitzero"`
	TotalPages   int `json:"totalPages,omitzero"`
	NextPage     int `json:"nextPage,omitzero"`
	PrevPage     int `json:"prevPage,omitzero"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, nil, map[string]string{"error": msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, info *Info, payload any) {
	var jsonData []byte
	var err error

	// Case 1: Error responses (>=400) - send raw payload
	switch {
	case code >= http.StatusBadRequest && info == nil:
		jsonData, err = json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		response := ResponseWrapper{
			Info: info,
			Data: payload,
		}
		jsonData, err = json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if jsonData == nil {
		jsonData = []byte("{}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
