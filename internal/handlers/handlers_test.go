package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ShvetsovYura/metrics-collector/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockStorage struct {
	Metrics map[string]types.Stringer
}

func (m *MockStorage) UpdateGauge(name string, val float64) {
	m.Metrics[name] = types.Gauge(val)
}

func (m *MockStorage) UpdateCounter(name string, val int64) {
	if v, ok := m.Metrics[name]; ok {
		m.Metrics[name] = v.(types.Counter) + types.Counter(val)
	} else {
		m.Metrics[name] = types.Counter(val)
	}
}

func (m *MockStorage) GetVal(name string) (types.Stringer, error) {
	if val, ok := m.Metrics[name]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MockStorage) ToList() []string {
	var list []string
	for _, c := range m.Metrics {
		list = append(list, c.ToString())
	}
	return list
}

type wantGauge struct {
	code  int
	mn    string
	val   types.Gauge
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
	m := new(MockStorage)
	m.Metrics = make(map[string]types.Stringer)
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
				val:   types.Gauge(3.4),
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
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := m.GetVal(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}
}

type wantCounter struct {
	code  int
	mn    string
	val   types.Counter
	isErr bool
}

func TestMetricUpdateCounterHandler(t *testing.T) {
	m := new(MockStorage)
	m.Metrics = make(map[string]types.Stringer)
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
				val:   types.Counter(3),
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
			assert.Equal(t, test.want.code, resp.StatusCode)
			if !test.want.isErr {
				v, err := m.GetVal(test.want.mn)
				require.Nil(t, err)
				assert.Equal(t, test.want.val.ToString(), v.ToString())
			}
		})
	}

}

func TestMetricGetValueHandler(t *testing.T) {
	m := new(MockStorage)
	m.Metrics = make(map[string]types.Stringer)
	m.Metrics["Alloc"] = types.Gauge(3.1234)
	m.Metrics["PullCounter"] = types.Counter(12345)
	m.Metrics["OtherMetric"] = types.Gauge(-123.30)

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
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}

func TestMetricGetAllValueHandler(t *testing.T) {
	m := new(MockStorage)
	m.Metrics = make(map[string]types.Stringer)
	m.Metrics["Alloc"] = types.Gauge(3.1234)
	m.Metrics["PullCounter"] = types.Counter(12345)
	m.Metrics["OtherMetric"] = types.Gauge(-123.30)

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
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
	}
}
