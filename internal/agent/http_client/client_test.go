package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestMetricHttpClient_Send(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := io.ReadAll(r.Body)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Log.Fatal(err.Error())
			}
		}()
		if len(req) != 69 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("OK")
		w.WriteHeader(http.StatusOK)
	})
	tu := httptest.NewServer(handler)
	defer tu.Close()

	dataToSend := agent.MetricItem{ID: "MemFree", MType: "gauge", Value: 123.456, Delta: 0}
	type fields struct {
		client        http.Client
		url           string
		contentType   string
		hashKey       string
		publicKeyPath string
	}
	type args struct {
		item agent.MetricItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantOut string
	}{
		{
			name: "simple request",
			fields: fields{
				client:        *tu.Client(),
				url:           tu.URL,
				contentType:   "application/json",
				hashKey:       "",
				publicKeyPath: "",
			},
			args: args{
				item: dataToSend,
			},
			wantErr: false,
			wantOut: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MetricHTTPClient{
				client:        tt.fields.client,
				url:           tt.fields.url,
				contentType:   tt.fields.contentType,
				hashKey:       tt.fields.hashKey,
				publicKeyPath: tt.fields.publicKeyPath,
			}
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			err := c.Send(tt.args.item, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricHttpClient.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = w.Close()
			assert.NoError(t, err)
			os.Stdout = old
			out, _ := io.ReadAll(r)
			got := string(out)
			assert.Equal(t, "OK\n", got)

		})
	}
}
