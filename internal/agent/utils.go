package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func sendMetric(data []byte, link string, contentType string, key string) error {
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
	if key != "" {
		hash := Hash(buf.Bytes(), key)
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

func Hash(value []byte, key string) string {
	h := sha256.New()
	h.Write(value)
	res := h.Sum(nil)
	return hex.EncodeToString(res)
}
