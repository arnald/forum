package helpers

import (
	"encoding/json"
	"net/http"
)

func ParseBodyRequest(r *http.Request, v any) (any, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
