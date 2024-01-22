package util

import (
	"io"
	"net/http"
)

func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func SendRequest(link string, contentType string, body io.Reader) error {
	r, err := http.Post(link, contentType, body) // "text/html"
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
