package util

import (
	"bytes"
	"net/http"

	"golang.org/x/exp/constraints"
)

func Contains(s []string, v string) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

func SendRequest(link string, contentType string, body []byte) error {

	reader := bytes.NewReader(body)
	resp, err := http.Post(link, contentType, reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func SaveOne[T constraints.Integer | constraints.Float](metricName string, metricValue T) {

}
