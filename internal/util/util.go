package util

import "net/http"

func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func SendRequest(link string) error {
	r, err := http.Post(link, "text/html", nil)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
