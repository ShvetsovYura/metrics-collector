package agent

import (
	"bytes"
	"compress/gzip"
	"net/http"
)

func sendMetric(data []byte, link string, contentType string) error {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)

	_, err := gzw.Write(data)
	if err != nil {
		return err
	}

	gzw.Close()

	req, err := http.NewRequest("POST", link, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", contentType)
	client := http.Client{}
	client.Do(req)
	return nil
}
