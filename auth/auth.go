package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API key from
// the headers of an http request
//
// Example:
//
// Authorization: ApiKey (insert apiKey here)
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authentication header is required")
	}

	vals := strings.Split(authHeader, " ")
	if len(vals) != 2 {
		return "", errors.New("Authentication header is invalid")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("Authentication header format is wrong")
	}

	return vals[1], nil
}
