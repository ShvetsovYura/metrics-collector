package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wantGauge struct {
	code  int
	mn    string
	val   metric.Gauge
	isErr bool
}

func (v *wantGauge) ToString() string {
	return strconv.FormatFloat(float64(v.val), 'f', -1, 64)

}
func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
func TestMetricUpdateGaugeHandler(t *testing.T) {
	m := memory.NewStorage(40)
	router := ServerRouter(m)
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		name   string
		path   string
		method string
		want   wantGauge
	}{
		{
			name:   "positive test #1",
			path:   "/update/gauge/Alloc/3.400",
			method: http.MethodPost,
			want: wantGauge{
				code:  http.StatusOK,
				val:   metric.Gauge(3.4),
				mn:    "Alloc",
				isErr: false,
			},
		},
		{
			name:   "not allowed",
			path:   "/update/gauge/Alloc/3.4",
			method: http.MethodGet,
			want: wantGauge{
				code:  http.StatusMethodNotAllowed,
				isErr: true,
			},
		}, {
			name:   "incorrect metric value",
			path:   "/update/gauge/Alloc/abracadabra",
			method: http.MethodPost,
			want: wantGauge{
				code:  http.StatusBadRequest,
				isErr: true,
			},
		},
		{
			name:   "incorrect metric type",
			path:   "/update/pipa/Alloc/123.23",
			method: http.MethodPost,
			want: wantGauge{
				code:  http.StatusBadRequest,
				isErr: true,
			},
		},
		{
			name:   "no metric value",
			path:   "/update/gauge/Alloc",
			method: http.MethodPost,
			want: wantGauge{
				code:  http.StatusNotFound,
				isErr: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, test.path)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := m.GetGauge(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}
}

type wantCounter struct {
	code  int
	mn    string
	val   metric.Counter
	isErr bool
}

func TestMetricUpdateCounterHandler(t *testing.T) {
	m := memory.NewStorage(40)
	router := ServerRouter(m)
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		name   string
		path   string
		method string
		want   wantCounter
	}{

		{
			name:   "positive test #2",
			path:   "/update/counter/PullCounter/3",
			method: http.MethodPost,
			want: wantCounter{
				code:  http.StatusOK,
				val:   metric.Counter(3),
				mn:    "PullCounter",
				isErr: false,
			},
		}, {
			name:   "wrong method ",
			path:   "/update/counter/PullCounter/332",
			method: http.MethodGet,
			want: wantCounter{
				code:  http.StatusMethodNotAllowed,
				isErr: true,
			},
		}, {
			name:   "wrong value ",
			path:   "/update/counter/PullCounter/332.234",
			method: http.MethodPost,
			want: wantCounter{
				code:  http.StatusBadRequest,
				isErr: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, test.path)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := m.GetGauge(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}

}

func TestMetricGetValueHandler(t *testing.T) {
	m := memory.NewStorage(40)
	m.UpdateGauge("Alloc", 3.1234)
	m.UpdateCounter(12345)
	m.UpdateGauge("OtherMetric", -123.30)

	router := ServerRouter(m)
	ts := httptest.NewServer(router)
	defer ts.Close()

	var testCases = []struct {
		url    string
		want   string
		status int
	}{
		{"/value/gauge/Alloc", "3.1234", http.StatusOK},
		{"/value/counter/PullCounter", "12345", http.StatusOK},
		{"/value/gauge/OtherMetric", "-123.3", http.StatusOK},
		{"/value/abra/Alloc", "", http.StatusNotFound},
		{"/value/gauge/abra", "", http.StatusNotFound},
	}

	for _, test := range testCases {
		resp, get := testRequest(t, ts, http.MethodGet, test.url)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}

func TestMetricGetAllValueHandler(t *testing.T) {
	m := memory.NewStorage(40)
	m.UpdateGauge("Alloc", 3.1234)
	m.UpdateCounter(12345)
	m.UpdateGauge("OtherMetric", -123.30)

	router := ServerRouter(m)
	ts := httptest.NewServer(router)
	defer ts.Close()

	var testCases = []struct {
		url    string
		want   string
		status int
	}{
		{"/", "3.1234, 12345, -123.3", http.StatusOK},
	}

	for _, test := range testCases {
		resp, get := testRequest(t, ts, http.MethodGet, test.url)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		for _, v := range strings.Split(test.want, ", ") {
			assert.Contains(t, get, v)
		}
	}
}
