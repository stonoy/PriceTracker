package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return "", errors.New("no authorization header was provided")
	}

	headerSlice := strings.Fields(header)

	// if headerSlice[0] != "Apikey" {
	// 	return "", errors.New("no authorization header name Apikey was provided")
	// }

	if len(headerSlice) < 2 {
		return "", errors.New("no token was provided")
	}

	return headerSlice[1], nil

}
