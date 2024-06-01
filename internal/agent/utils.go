package agent

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func sendMetric(data []byte, link string, contentType string, key string) error {
	var buf bytes.Buffer

	req, err := http.NewRequest("POST", link, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentType)

	if util.Contains([]string{"application/json", "text/html"}, contentType) {
		req.Header.Add("Content-Encoding", "gzip")
		req.Header.Add("Accept-Encoding", "gzip")
		gzw := gzip.NewWriter(&buf)

		_, err := gzw.Write(data)
		if err != nil {
			return err
		}

		gzw.Close()
	} else {
		writer := io.Writer(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return err
		}
	}

	if key != "" {
		hash := util.Hash(buf.Bytes(), key)
		req.Header.Add("HashSHA256", hash)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
