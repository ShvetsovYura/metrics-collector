package httpclient

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

type MetricHTTPClient struct {
	client        http.Client
	url           string
	contentType   string
	hashKey       string
	publicKeyPath string
}

func NewMetricSender(url string, contentType string, hashKey string, publicKeyPath string) *MetricHTTPClient {
	return &MetricHTTPClient{
		client:        http.Client{},
		url:           url,
		contentType:   contentType,
		hashKey:       hashKey,
		publicKeyPath: publicKeyPath,
	}
}

func (c *MetricHTTPClient) Send(data []byte) error {
	var buf bytes.Buffer
	var headers = http.Header{}
	var data_ []byte

	headers.Add("Content-Type", c.contentType)
	addresses, err := util.GetLocalIPs()
	if err != nil {
		return fmt.Errorf("ошибка получения IP %w", err)
	}
	headers.Add("X-Real-IP", addresses[0].String())

	if c.publicKeyPath != "" {
		var errEncrypt error
		data_, errEncrypt = util.EncryptData(data, c.publicKeyPath)
		if errEncrypt != nil {
			return fmt.Errorf("ошибка при шифровании сообщения %w", errEncrypt)
		}
	} else {
		data_ = data
	}

	if util.Contains([]string{"application/json", "text/html"}, c.contentType) {
		headers.Add("Content-Encoding", "gzip")
		headers.Add("Accept-Encoding", "gzip")

		gzw := gzip.NewWriter(&buf)

		_, err := gzw.Write(data_)
		if err != nil {
			return fmt.Errorf("ошибка при записи gzip тела при отправке, %w", err)
		}

		err = gzw.Close()
		if err != nil {
			return fmt.Errorf("ошибка при закрытии gzip писателя, %w", err)
		}
	} else {
		writer := io.Writer(&buf)

		_, err := writer.Write(data_)
		if err != nil {
			return fmt.Errorf("ошибка записи тела web запроса, %w", err)
		}
	}

	if c.hashKey != "" {
		hash := util.Hash(buf.Bytes(), c.hashKey)
		headers.Add("HashSHA256", hash)
	}

	req, err := http.NewRequest("POST", c.url, &buf)
	if err != nil {
		return fmt.Errorf("ошибка создания web запроса для отправки метрик, %w", err)
	}
	req.Header = headers
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения web запроса, %w", err)
	}

	logger.Log.Infof("response status code: %d", resp.StatusCode)

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Printf("error on close response body, %s", err.Error())
		}
	}()

	return nil
}
