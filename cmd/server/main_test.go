package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHandlere(t *testing.T) {
	type want struct {
		code int
		// response    string
		// contentType string
	}

	tests := []struct {
		name   string
		path   string
		method string
		want   want
	}{
		{
			name:   "positive test #1",
			path:   "/update/gauge/Alloc/3.4",
			method: http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		}, {
			name:   "positive test #2",
			path:   "/update/counter/PullCounter/3",
			method: http.MethodPost,
			want: want{
				code: http.StatusOK,
			},
		}, {
			name:   "not allowed",
			path:   "/update/gauge/Alloc/3.4",
			method: http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		}, {
			name:   "incorrect metric value",
			path:   "/update/gauge/Alloc/abracadabra",
			method: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "incorrect metric type",
			path:   "/update/pipa/Alloc/123.23",
			method: http.MethodPost,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:   "no metric value",
			path:   "/update/gauge/Alloc",
			method: http.MethodPost,
			want: want{
				code: http.StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			metricHandler(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
