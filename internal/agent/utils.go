package agent

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func sendMetric(data []byte, link string, contentType string) error {
	var buf bytes.Buffer
	var writer io.Writer
	if util.Contains([]string{"application/json", "text/html"}, contentType) {
		gzw := gzip.NewWriter(&buf)

		_, err := gzw.Write(data)
		if err != nil {
			return err
		}

		gzw.Close()

	} else {
		writer = io.Writer(&buf)
		writer.Write(data)
	}

	req, err := http.NewRequest("POST", link, &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept-Encoding", "gzip")
	client := http.Client{}
	resp, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}
	defer resp.Body.Close()
	return nil
}
