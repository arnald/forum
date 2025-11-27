package helpers

import (
	"encoding/json"
	"net/http"
)

// DecodeBackendResponse is a generic helper to decode wrapped backend responses.
// It expects responses in the format: { "data": { ...fields... } }.
func DecodeBackendResponse[T any](resp *http.Response, target *T) error {
	wrapper := struct {
		Data T `json:"data"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return err
	}

	*target = wrapper.Data

	return nil
}
