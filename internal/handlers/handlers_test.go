package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
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
func testRequest(t *testing.T, ts *httptest.Server, method, path string, data []byte) (*http.Response, string) {
	var buf *bytes.Buffer
	var req *http.Request
	var reqErr error
	if len(data) > 0 {
		buf = bytes.NewBuffer(data)
		req, reqErr = http.NewRequest(method, ts.URL+path, buf)
	} else {
		req, reqErr = http.NewRequest(method, ts.URL+path, nil)

	}

	require.NoError(t, reqErr)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricSetGaugeHandler(t *testing.T) {
	fs := file.NewFileStorage("tt.txt", 40, false)
	router := ServerRouter(fs)
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
			resp, _ := testRequest(t, ts, test.method, test.path, []byte{})
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := fs.GetGauge(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}
}

func ServerRouter(fs *file.FileStorage) {
	panic("unimplemented")
}

type wantCounter struct {
	code  int
	mn    string
	val   metric.Counter
	isErr bool
}

func TestMetricSetCounterHandler(t *testing.T) {
	m := memory.NewMemStorage(40)
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
			name:   "positive test #1",
			path:   "/update/counter/PollCount/3",
			method: http.MethodPost,
			want: wantCounter{
				code:  http.StatusOK,
				val:   metric.Counter(3),
				mn:    "PollCount",
				isErr: false,
			},
		}, {
			name:   "wrong method ",
			path:   "/update/counter/PollCount/332",
			method: http.MethodGet,
			want: wantCounter{
				code:  http.StatusMethodNotAllowed,
				isErr: true,
			},
		}, {
			name:   "wrong value ",
			path:   "/update/counter/PollCount/332.234",
			method: http.MethodPost,
			want: wantCounter{
				code:  http.StatusBadRequest,
				isErr: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.method, test.path, nil)
			defer resp.Body.Close()
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := m.GetCounter(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}

}

func TestMetricGetValueHandler(t *testing.T) {
	m := memory.NewMemStorage(40)
	m.SetGauge("Alloc", 3.1234)
	m.SetCounter("PollCount", 12345)
	m.SetGauge("OtherMetric", -123.30)

	router := ServerRouter(m)
	ts := httptest.NewServer(router)
	defer ts.Close()

	var testCases = []struct {
		url    string
		want   string
		status int
	}{
		{"/value/gauge/Alloc", "3.1234", http.StatusOK},
		{"/value/counter/PollCount", "12345", http.StatusOK},
		{"/value/gauge/OtherMetric", "-123.3", http.StatusOK},
		{"/value/abra/Alloc", "", http.StatusNotFound},
		{"/value/gauge/abra", "", http.StatusNotFound},
	}

	for _, test := range testCases {
		resp, get := testRequest(t, ts, http.MethodGet, test.url, nil)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}

func TestMetricGetAllValueHandler1(t *testing.T) {
	m := memory.NewMemStorage(40)
	m.SetGauge("Alloc", 3.1234)
	m.SetCounter("PollCount", 12345)
	m.SetGauge("OtherMetric", -123.30)

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
		resp, get := testRequest(t, ts, http.MethodGet, test.url, nil)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		for _, v := range strings.Split(test.want, ", ") {
			assert.Contains(t, get, v)
		}
	}
}

func TestMetricUpdateHandler(t *testing.T) {
	countMetrics := 40
	fsPath := "/tmp/myFileStorage.txt"
	fs := file.NewFileStorage(fsPath, countMetrics, true)

	router := ServerRouter(fs)
	ts := httptest.NewServer(router)
	defer func() {
		ts.Close()
		os.Remove(fsPath)
	}()
	wantValues := []float64{3.400}
	wantDeltas := []int64{3}
	tests := []struct {
		name       string
		reqData    string
		want       models.Metrics
		wantStatus int
	}{

		{
			name:    "positive update counter",
			reqData: `{"id":"PollCounter", "type":"counter", "delta": 3}`,

			want:       models.Metrics{ID: "PollCounter", Delta: &wantDeltas[0], MType: internal.InCounterName},
			wantStatus: http.StatusOK,
		}, {
			name:       "wrong value ",
			reqData:    `{"id":"PollCounter", "type":"counter", "delta": 332.234}`,
			want:       models.Metrics{},
			wantStatus: http.StatusBadRequest,
		}, {
			name:       "positive update gauge",
			reqData:    `{"id":"Alloc", "type":"gauge", "value": 3.400}`,
			want:       models.Metrics{ID: "Alloc", Value: &wantValues[0], MType: internal.InGaugeName},
			wantStatus: http.StatusOK,
		},
		{
			name:       "incorrect metric value",
			reqData:    `{"id":"Alloc", "type":"gauge", "value": "abracadabra"}`,
			want:       models.Metrics{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "incorrect metric type",
			reqData:    `{"id":"Alloc", "type":"other", "value": 123.45}`,
			want:       models.Metrics{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, "POST", "/update/", []byte(test.reqData))
			defer resp.Body.Close()
			assert.Equal(t, test.wantStatus, resp.StatusCode)
			if test.wantStatus == http.StatusOK {
				respJSON := &models.Metrics{}
				json.Unmarshal([]byte(body), respJSON)
				assert.Equal(t, test.want, *respJSON)
			}
		})
	}

}

func TestMetricValueHandler(t *testing.T) {
	countMetrics := 40
	fsPath := "/tmp/myFileStorage.txt"
	fs := file.NewFileStorage(fsPath, countMetrics, true)

	fs.SetGauge("Alloc", 3.1234)
	fs.SetCounter("PollCount", 12345)
	fs.SetGauge("OtherMetric", -123.30)

	router := ServerRouter(fs)
	ts := httptest.NewServer(router)
	defer func() {
		ts.Close()
		os.Remove(fsPath)
	}()

	wantValues := []float64{3.1234, -123.3}
	var wantDelta int64 = 12345
	var testCases = []struct {
		name    string
		reqData string
		want    models.Metrics
		status  int
	}{
		{
			name:    "check get correct gauge",
			reqData: `{"id":"Alloc", "type":"gauge"}`,
			want:    models.Metrics{ID: "Alloc", Value: &wantValues[0], MType: internal.InGaugeName},
			status:  http.StatusOK,
		},
		{
			name:    "check get other correct gauge",
			reqData: `{"id":"OtherMetric", "type":"gauge"}`,
			want:    models.Metrics{ID: "OtherMetric", Value: &wantValues[1], MType: internal.InGaugeName},
			status:  http.StatusOK,
		},
		{
			name:    "check get unknown gauge",
			reqData: `{"id":"Ugu", "type":"gauge"}`,
			want:    models.Metrics{},
			status:  http.StatusNotFound,
		},
		{
			name:    "check get counter",
			reqData: `{"id":"PollCount", "type":"counter"}`,
			want:    models.Metrics{ID: "PollCount", Delta: &wantDelta, MType: internal.InCounterName},
			status:  http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			b := []byte(test.reqData)
			resp, body := testRequest(t, ts, http.MethodPost, "/value/", b)
			defer resp.Body.Close()
			assert.Equal(t, test.status, resp.StatusCode)
			if test.status == http.StatusOK {

				respJSON := &models.Metrics{}
				json.Unmarshal([]byte(body), respJSON)
				assert.Equal(t, test.want, *respJSON)
			}
		})

	}
}

func TestMetricGetAllValueHandler(t *testing.T) {
	countMetrics := 40
	fsPath := "/tmp/myFileStorage.txt"
	fs := file.NewFileStorage(fsPath, countMetrics, true)

	fs.SetGauge("Alloc", 3.1234)
	fs.SetCounter("PollCount", 12345)
	fs.SetGauge("OtherMetric", -123.30)

	router := ServerRouter(fs)
	ts := httptest.NewServer(router)
	defer func() {
		ts.Close()
		os.Remove(fsPath)
	}()

	var testCases = []struct {
		url    string
		want   string
		status int
	}{
		{"/", "3.1234, 12345, -123.3", http.StatusOK},
	}

	for _, test := range testCases {
		resp, get := testRequest(t, ts, http.MethodGet, test.url, nil)
		defer resp.Body.Close()
		assert.Equal(t, test.status, resp.StatusCode)
		for _, v := range strings.Split(test.want, ", ") {
			assert.Contains(t, get, v)
		}
	}
}
