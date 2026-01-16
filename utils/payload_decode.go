package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)


func DecodeData (r *http.Request, payload any) error  {

	if r.Body == nil {
		return errors.New("Failed to load body and decode it!")
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return errors.New("Failed to decode data payload into a json type!")
	}

	return nil

}