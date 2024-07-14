package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func sendMetric(data []byte, link string, contentType string, key string) error {
	var buf bytes.Buffer

	req, err := http.NewRequest("POST", link, &buf)
	if err != nil {
		return fmt.Errorf("ошибка создания web запроса для отправки метрик, %w", err)
	}

	req.Header.Add("Content-Type", contentType)

	if util.Contains([]string{"application/json", "text/html"}, contentType) {
		req.Header.Add("Content-Encoding", "gzip")
		req.Header.Add("Accept-Encoding", "gzip")

		gzw := gzip.NewWriter(&buf)

		_, writeErr := gzw.Write(data)
		if writeErr != nil {
			return fmt.Errorf("ошибка при записи gzip тела при отправке, %w", writeErr)
		}

		err = gzw.Close()
		if err != nil {
			return fmt.Errorf("ошибка при закрытии gzip писателя, %w", err)
		}
	} else {
		writer := io.Writer(&buf)

		_, writeErr := writer.Write(data)
		if writeErr != nil {
			return fmt.Errorf("ошибка записи тела web запроса, %w", writeErr)
		}
	}

	if key != "" {
		hash := util.Hash(buf.Bytes(), key)
		req.Header.Add("HashSHA256", hash)
	}

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("ошибка выполнения web запроса, %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("error on close response body, %s", err.Error())
		}
	}()

	return nil
}
