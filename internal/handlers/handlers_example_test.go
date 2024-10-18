package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
)

func ExampleDBPingHandler() {
	s := storage.NewMemory(10)
	routes := handlers.ServerRouter(s, "abc", "", "")
	ts := httptest.NewServer(routes)

	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	fmt.Println(resp.Status)

	// Output:
	// 200 OK

}

func ExampleMetricBatchUpdateHandler() {
	s := storage.NewMemory(10)
	routes := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(routes)

	defer ts.Close()

	gaugeValue1 := models.Gauge(984.723)
	gaugeValue2 := models.Gauge(-234433.33)
	counterValue := models.Counter(4)
	metrics := []models.MetricItem{{
		ID:    "gaugeMetric1",
		MType: "gauge",
		Value: gaugeValue1.GetRawValue(),
	}, {
		ID:    "gaugeMetric1",
		MType: "gauge",
		Value: gaugeValue2.GetRawValue(),
	}, {
		ID:    "counterMetric",
		MType: "counter",
		Delta: counterValue.GetRawValue(),
	}}

	var body bytes.Buffer
	jsonEncoder := json.NewEncoder(&body)

	err := jsonEncoder.Encode(metrics)
	if err != nil {
		log.Fatalf("не удалось преобразовать в json, %s", err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/updates/", &body)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	fmt.Println(resp.Status)

	// Output:
	// 200 OK
}

func ExampleMetricGetCurrentValuesHandler() {
	s := storage.NewMemory(10)
	ctx := context.Background()
	s.SetGauges(ctx, map[string]float64{
		"memTotal":  345.43,
		"freeSpace": 9563738.322,
		"maxLoad":   97.34,
	})

	_ = s.SetCounter(ctx, "count", 4)

	routes := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(routes)

	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	// Output:
	// 9563738.322, 97.34, 345.43, 4
}

func ExampleMetricGetValueHandlerWithBody() {
	s := storage.NewMemory(10)
	ctx := context.Background()
	gauges := map[string]float64{
		"memTotal":  345.43,
		"freeSpace": 9563738.322,
		"maxLoad":   97.34,
	}
	s.SetGauges(ctx, gauges)
	routes := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(routes)

	defer ts.Close()

	reqBytes, _ := json.Marshal(models.MetricItem{
		ID: "memTotal", MType: "gauge",
	})
	reqBody := bytes.NewBuffer(reqBytes)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/value/", reqBody)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	var m models.MetricItem

	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &m)

	fmt.Println(resp.Status)

	fmt.Println(string(respBytes))
	// Output:
	// 200 OK
	// {"id":"memTotal","type":"gauge","value":345.43}

}

func ExampleMetricUpdateHandlerWithBody() {
	s := storage.NewMemory(10)
	r := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(r)

	defer ts.Close()

	gauge := models.Gauge(4011.1)
	reqBytes, _ := json.Marshal(models.MetricItem{
		ID:    "allocateMem",
		MType: "gauge",
		Value: gauge.GetRawValue(),
	})
	reqBuf := bytes.NewBuffer(reqBytes)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/update/", reqBuf)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	respBytes, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Status)
	fmt.Println(string(respBytes))

	// Output:
	// 200 OK
	// {"id":"allocateMem","type":"gauge","value":4011.1}
}

func ExampleMetricGetValueHandler() {
	s := storage.NewMemory(10)
	ctx := context.Background()
	gauges := map[string]float64{
		"allogMem":  3718.23,
		"freeMem":   1528.30,
		"usedSpace": 134672046.234,
	}
	s.SetGauges(ctx, gauges)
	r := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(r)

	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/value/gauge/freeMem", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	respBytes, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(respBytes))

	// Output:
	// 200 OK
	// 1528.3
}

func ExampleMetricUpdateHandler() {
	s := storage.NewMemory(10)
	r := handlers.ServerRouter(s, "key", "", "")
	ts := httptest.NewServer(r)

	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/update/gauge/allocMem/2139.43", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("ошибка при закрытии тела запроса, %s", err.Error())
		}
	}()

	fmt.Println(resp.Status)
	// Output:
	// 200 OK
}
