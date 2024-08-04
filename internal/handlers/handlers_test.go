package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
)

type wantGauge struct {
	code  int
	mn    string
	val   models.Gauge
	isErr bool
}

func (v *wantGauge) ToString() string {
	return strconv.FormatFloat(float64(v.val), 'f', -1, 64)
}
func testRequest(t *testing.T, ts *httptest.Server, method, path string, data []byte) (*http.Response, string) {
	var (
		buf    *bytes.Buffer
		req    *http.Request
		reqErr error
	)

	if len(data) > 0 {
		buf = bytes.NewBuffer(data)
		req, reqErr = http.NewRequest(method, ts.URL+path, buf)
	} else {
		req, reqErr = http.NewRequest(method, ts.URL+path, nil)
	}

	require.NoError(t, reqErr)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer func() {
		closeErr := resp.Body.Close()
		if err != nil {
			fmt.Printf("не удалось закрыть тело ответа, %s", closeErr.Error())
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricSetGaugeHandler(t *testing.T) {
	mem := storage.NewMemory(40)
	fs := storage.NewFile("tt.txt", mem, false, 0)
	router := ServerRouter(fs, "", "")
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
				val:   models.Gauge(3.4),
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

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, resp.StatusCode)

			if !test.want.isErr {
				v, err := fs.GetGauge(context.Background(), test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}
}

type wantCounter struct {
	code  int
	mn    string
	val   models.Counter
	isErr bool
}

func TestMetricSetCounterHandler(t *testing.T) {
	m := storage.NewMemory(40)
	router := ServerRouter(m, "", "")
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
				val:   models.Counter(3),
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

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
				}
			}()

			assert.Equal(t, test.want.code, resp.StatusCode)

			if !test.want.isErr {
				v, err := m.GetCounter(context.Background(), test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}
}

func TestMetricGetValueHandler(t *testing.T) {
	ctx := context.Background()
	m := storage.NewMemory(40)

	err := m.SetGauge(ctx, "Alloc", 3.1234)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = m.SetCounter(ctx, "PollCount", 12345)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = m.SetGauge(ctx, "OtherMetric", -123.30)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	router := ServerRouter(m, "", "")
	ts := httptest.NewServer(router)

	defer ts.Close()

	var tests = []struct {
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

	for _, test := range tests {
		resp, get := testRequest(t, ts, http.MethodGet, test.url, nil)
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
			}
		}()
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}

func TestMetricGetAllValueHandler1(t *testing.T) {
	m := storage.NewMemory(40)
	ctx := context.Background()

	err := m.SetGauge(ctx, "Alloc", 3.1234)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = m.SetCounter(ctx, "PollCount", 12345)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = m.SetGauge(ctx, "OtherMetric", -123.30)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	router := ServerRouter(m, "", "")
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
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
			}
		}()
		assert.Equal(t, test.status, resp.StatusCode)

		for _, v := range strings.Split(test.want, ", ") {
			assert.Contains(t, get, v)
		}
	}
}

func TestMetricUpdateHandler(t *testing.T) {
	mem := storage.NewMemory(40)
	fsPath := "/tmp/myFileStorage.txt"
	fs := storage.NewFile(fsPath, mem, true, 0)

	router := ServerRouter(fs, "", "")
	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()

		err := os.Remove(fsPath)
		if err != nil {
			fmt.Printf("ошибка удаления файла, %s", err.Error())
		}
	}()

	wantValues := []float64{3.400}
	wantDeltas := []int64{3}
	tests := []struct {
		name       string
		reqData    string
		want       models.MetricItem
		wantStatus int
	}{

		{
			name:    "positive update counter",
			reqData: `{"id":"PollCounter", "type":"counter", "delta": 3}`,

			want:       models.MetricItem{ID: "PollCounter", Delta: &wantDeltas[0], MType: internal.InCounterName},
			wantStatus: http.StatusOK,
		}, {
			name:       "wrong value ",
			reqData:    `{"id":"PollCounter", "type":"counter", "delta": 332.234}`,
			want:       models.MetricItem{},
			wantStatus: http.StatusBadRequest,
		}, {
			name:       "positive update gauge",
			reqData:    `{"id":"Alloc", "type":"gauge", "value": 3.400}`,
			want:       models.MetricItem{ID: "Alloc", Value: &wantValues[0], MType: internal.InGaugeName},
			wantStatus: http.StatusOK,
		},
		{
			name:       "incorrect metric value",
			reqData:    `{"id":"Alloc", "type":"gauge", "value": "abracadabra"}`,
			want:       models.MetricItem{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "incorrect metric type",
			reqData:    `{"id":"Alloc", "type":"other", "value": 123.45}`,
			want:       models.MetricItem{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, "POST", "/update/", []byte(test.reqData))
			defer func() {
				err := resp.Body.Close()
				if err != nil {
					t.Fatalf("не удалось закрыть тело ответа, %s", err.Error())
				}
			}()
			assert.Equal(t, test.wantStatus, resp.StatusCode)

			if test.wantStatus == http.StatusOK {
				respJSON := &models.MetricItem{}

				err := json.Unmarshal([]byte(body), respJSON)
				if err != nil {
					t.Fatalf("ошибка преобразования в объект, %s", err.Error())
				}

				assert.Equal(t, test.want, *respJSON)
			}
		})
	}
}

func TestMetricValueHandler(t *testing.T) {
	mem := storage.NewMemory(40)
	fsPath := "/tmp/myFileStorage.txt"
	fs := storage.NewFile(fsPath, mem, true, 0)
	ctx := context.Background()

	err := fs.SetGauge(ctx, "Alloc", 3.1234)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = fs.SetCounter(ctx, "PollCount", 12345)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = fs.SetGauge(ctx, "OtherMetric", -123.30)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	router := ServerRouter(fs, "", "")
	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()

		err := os.Remove(fsPath)
		if err != nil {
			log.Fatalf("ошибка удаления файла, %s", err.Error())
		}
	}()

	wantValues := []float64{3.1234, -123.3}

	var wantDelta int64 = 12345

	var testCases = []struct {
		name    string
		reqData string
		want    models.MetricItem
		status  int
	}{
		{
			name:    "check get correct gauge",
			reqData: `{"id":"Alloc", "type":"gauge"}`,
			want:    models.MetricItem{ID: "Alloc", Value: &wantValues[0], MType: internal.InGaugeName},
			status:  http.StatusOK,
		},
		{
			name:    "check get other correct gauge",
			reqData: `{"id":"OtherMetric", "type":"gauge"}`,
			want:    models.MetricItem{ID: "OtherMetric", Value: &wantValues[1], MType: internal.InGaugeName},
			status:  http.StatusOK,
		},
		{
			name:    "check get unknown gauge",
			reqData: `{"id":"Ugu", "type":"gauge"}`,
			want:    models.MetricItem{},
			status:  http.StatusNotFound,
		},
		{
			name:    "check get counter",
			reqData: `{"id":"PollCount", "type":"counter"}`,
			want:    models.MetricItem{ID: "PollCount", Delta: &wantDelta, MType: internal.InCounterName},
			status:  http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			b := []byte(test.reqData)
			resp, body := testRequest(t, ts, http.MethodPost, "/value/", b)

			defer func() {
				err := resp.Body.Close()
				if err != nil {
					fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
				}
			}()

			assert.Equal(t, test.status, resp.StatusCode)

			if test.status == http.StatusOK {
				respJSON := &models.MetricItem{}
				_ = json.Unmarshal([]byte(body), respJSON)
				assert.Equal(t, test.want, *respJSON)
			}
		})
	}
}

func TestMetricGetAllValueHandler(t *testing.T) {
	mem := storage.NewMemory(40)
	fsPath := "/tmp/myFileStorage.txt"
	fs := storage.NewFile(fsPath, mem, true, 0)
	ctx := context.Background()

	err := fs.SetGauge(ctx, "Alloc", 3.1234)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = fs.SetCounter(ctx, "PollCount", 12345)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	err = fs.SetGauge(ctx, "OtherMetric", -123.30)
	if err != nil {
		t.Fatalf("не удалось установить метрику, %s", err.Error())
	}

	router := ServerRouter(fs, "", "")
	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()

		err := os.Remove(fsPath)
		if err != nil {
			log.Fatalf("ошибка удаления файла, %s", err.Error())
		}
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

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
			}
		}()

		assert.Equal(t, test.status, resp.StatusCode)

		for _, v := range strings.Split(test.want, ", ") {
			assert.Contains(t, get, v)
		}
	}
}

func TestMetricBatchUpdateHandler(t *testing.T) {
	mem := storage.NewMemory(40)
	router := ServerRouter(mem, "", "")
	ts := httptest.NewServer(router)

	defer func() {
		ts.Close()
	}()

	gaugeWants := []models.Gauge{123.56, 0.0}
	counterWants := []models.Counter{0, 112, 1}
	tests := []struct {
		name       string
		input      []models.MetricItem
		want       []models.Gauge
		wantStatus int
	}{{
		name: "many gauge metrics",
		input: []models.MetricItem{
			{
				ID:    "metric1",
				MType: "gauge",
				Value: gaugeWants[0].GetRawValue(),
			}, {
				ID:    "metric2",
				MType: "gauge",
				Value: gaugeWants[1].GetRawValue(),
			},
		},
		want:       nil,
		wantStatus: 200,
	}, {
		name: "mnny counter metirc",
		input: []models.MetricItem{
			{
				ID:    "counter_metric1",
				MType: "counter",
				Delta: counterWants[0].GetRawValue(),
			}, {
				ID:    "counter_metric2",
				MType: "counter",
				Delta: counterWants[1].GetRawValue(),
			}, {
				ID:    "counter_metric3",
				MType: "counter",
				Delta: counterWants[2].GetRawValue(),
			},
		},
		want:       nil,
		wantStatus: 200,
	},
	}

	for _, test := range tests {
		reqData, _ := json.Marshal(test.input)
		resp, _ := testRequest(t, ts, http.MethodPost, "/updates/", reqData)

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Printf("не удалось закрыть тело ответа, %s", err.Error())
			}
		}()

		assert.Equal(t, test.wantStatus, resp.StatusCode)
	}
}
