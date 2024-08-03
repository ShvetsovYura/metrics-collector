package middlewares

import (
	"bufio"
	"fmt"
	"io"

	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

type decryptReader struct {
	req            io.ReadCloser
	reader         io.Reader
	privateKeyPath string
}

func newDecryptReader(r io.ReadCloser, privateKeyPath string) *decryptReader {
	return &decryptReader{
		req:            r,
		reader:         bufio.NewReader(r),
		privateKeyPath: privateKeyPath,
	}
}

func (dr *decryptReader) Read(p []byte) (n int, err error) {
	decrytedMessage, err := util.DecryptData(p, dr.privateKeyPath)
	if err != nil {
		return 0, fmt.Errorf("error on decrypt message %e", err)
	}
	return dr.reader.Read(decrytedMessage)
}
func (dr *decryptReader) Close() error {
	if err := dr.req.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия decryptReader, %w", err)
	}

	return nil
}

func DecryptMessage(privateKeyPath string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		handler := func(w http.ResponseWriter, req *http.Request) {
			if privateKeyPath != "" {
				decrReader := newDecryptReader(req.Body, privateKeyPath)

				req.Body = decrReader
				defer func() {
					errClose := decrReader.Close()
					if errClose != nil {
						logger.Log.Errorf("ошибка закрытия decryptReader, %e", errClose)
					}
				}()
			}
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(handler)
	}
}
